package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	serverConfig "github.com/ramil063/gometrics/cmd/server/config"
	"github.com/ramil063/gometrics/cmd/server/handlers"
	serverGRPC "github.com/ramil063/gometrics/cmd/server/handlers/grpc/server"
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/cmd/server/storage/file"
	"github.com/ramil063/gometrics/internal/constants"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}

	params := serverConfig.NewConfigParams(
		constants.ConfigHTTPConsoleShortKey,
		constants.ConfigHTTPConsoleFullKey,
		constants.ConfigHTTPTypeAlias)
	config, err := serverConfig.GetConfig(params)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "config")
	}
	handlers.InitFlags(config)

	manager := crypto.NewCryptoManager()
	if handlers.CryptoKey != "" {
		defaultDecryptor, decryptorErr := crypto.NewRSADecryptor(handlers.CryptoKey)
		if decryptorErr != nil {
			logger.WriteErrorLog(decryptorErr.Error(), "Failed to create decryptor")
		}
		manager.SetDefaultDecryptor(defaultDecryptor)
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	var s = server.GetStorage(handlers.FileStoragePath, handlers.DatabaseDSN)

	if handlers.DatabaseDSN != "" {
		rep, errRepo := dml.NewRepository()
		if errRepo != nil {
			logger.WriteErrorLog(errRepo.Error(), "NewRepository")
			return
		}
		dml.DBRepository = *rep
		err = db.Init(&dml.DBRepository)
		defer dml.DBRepository.Close()
		if err != nil {
			logger.WriteErrorLog("Error in init db", err.Error())
			return
		}
	}

	writingToFileIsEnabledAndAvailable := handlers.FileStoragePath != ""
	if handlers.StoreInterval > 0 && writingToFileIsEnabledAndAvailable {
		if !handlers.Restore {
			// работаем с новыми метриками, очищая файл со старыми
			err = file.ClearFileContent(handlers.FileStoragePath)
			if err != nil {
				logger.WriteErrorLog(err.Error(), "ClearFileContent")
			}
		}
		ticker := time.NewTicker(time.Duration(handlers.StoreInterval) * time.Second)
		go func() {
			err = server.SaveMetricsPerTime(server.MaxSaverWorkTime, ticker, s)
			if err != nil {
				logger.WriteErrorLog(err.Error(), "SaveMetricsPerTime")
			}
		}()
	}

	srv := &http.Server{
		Addr:    handlers.MainURL,
		Handler: server.Router(s, manager),
	}

	grpcFlags, grpcStorage, manager, err := serverGRPC.PrepareServerEnvironment()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "serverGRPC.PrepareServerEnvironment")
	}

	grpcServer, err := serverGRPC.GetGRPCServer(grpcFlags, grpcStorage, manager)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "GetGRPCServer init error")
	}

	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})
	// регистрируем перенаправление прерываний
	ctxGrSh, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// запускаем горутину обработки пойманных прерываний
	go func() {
		<-ctxGrSh.Done()
		log.Println("Starting graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// получили сигнал запускаем процедуру graceful shutdown
		if err = srv.Shutdown(ctx); err != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown error: %v", err)
		}
		grpcServer.GracefulStop()

		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
		log.Println("All connections closed")
	}()

	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Printf("HTTP server ListenAndServe error: %v", err)
		stop() // Триггерим shutdown при ошибке сервера
	}
	// ждём завершения процедуры graceful shutdown
	<-idleConnsClosed
	// получили оповещение о завершении
	// здесь можно освобождать ресурсы перед выходом,
	// например закрыть соединение с базой данных,
	// закрыть открытые файлы
	// если для хранения использовались
	// 1) база данных
	// Явное освобождение ресурса
	if err = dml.DBRepository.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}
	// 2) файловая система - то файл закроется,
	//    так как везде при открытии файла рядом есть defer Reader.Close() или defer Writer.Close()
	// 3) оперативная память - то все мьютексы освободятся,
	//    так как везде при доступе к мапе есть defer ms.mx.Unlock() и defer ms.mx.RUnlock()
	fmt.Println("Server Shutdown gracefully")
}
