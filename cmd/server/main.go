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
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/cmd/server/storage/file"
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

	config, err := serverConfig.GetConfig()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "config")
	}
	handlers.InitFlags(config)

	if handlers.CryptoKey != "" {
		crypto.DefaultDecryptor, err = crypto.NewRSADecryptor(handlers.CryptoKey)
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Failed to create decryptor")
		}
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
		Handler: server.Router(s),
	}

	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})
	// канал для перенаправления прерываний
	// поскольку нужно отловить всего одно прерывание,
	// ёмкости 1 для канала будет достаточно
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// запускаем горутину обработки пойманных прерываний
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err = srv.Shutdown(ctx); err != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown error: %v", err)
		}

		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	// ждём завершения процедуры graceful shutdown
	<-idleConnsClosed
	// получили оповещение о завершении
	// здесь можно освобождать ресурсы перед выходом,
	// например закрыть соединение с базой данных,
	// закрыть открытые файлы
	// если для хранения использовались
	// 1) база данных - то она закроется,
	//    так как сразу после открытия базы есть defer dml.DBRepository.Close()
	// 2) файловая система - то файл закроется,
	//    так как везде при открытии файла рядом есть defer Reader.Close() или defer Writer.Close()
	// 3) оперативная память - то все мьютексы освободятся,
	//    так как везде при доступе к мапе есть defer ms.mx.Unlock() и defer ms.mx.RUnlock()
	fmt.Println("Server Shutdown gracefully")
}
