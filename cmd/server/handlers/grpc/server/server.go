package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	serverConfig "github.com/ramil063/gometrics/cmd/server/config"
	grpcHandlers "github.com/ramil063/gometrics/cmd/server/handlers/grpc"
	"github.com/ramil063/gometrics/cmd/server/handlers/grpc/interceptors"
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/cmd/server/storage/file"
	"github.com/ramil063/gometrics/internal/constants"
	pb "github.com/ramil063/gometrics/internal/grpc/proto"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

// PrepareServerEnvironment подготавливает окружение для работы сервера
func PrepareServerEnvironment() (*grpcHandlers.ServerConfigFlags, server.Storager, *crypto.Manager, error) {
	paramsGRPC := serverConfig.NewConfigParams(
		constants.ConfigGRPCConsoleShortKey,
		constants.ConfigGRPCConsoleFullKey,
		constants.ConfigGRPCTypeAlias)

	configGRPC, err := serverConfig.GetConfig(paramsGRPC)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "GetConfig")
	}

	flagsGRPC, err := grpcHandlers.GetFlags(configGRPC)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "GetFlags")
		return nil, nil, nil, err
	}

	var grpcStorage = server.GetStorage(flagsGRPC.FileStoragePath, flagsGRPC.DatabaseDSN)

	if flagsGRPC.DatabaseDSN != "" {
		rep, errRepo := dml.NewRepository()
		if errRepo != nil {
			logger.WriteErrorLog(errRepo.Error(), "NewRepository")
			return nil, nil, nil, errRepo
		}
		dml.DBRepository = *rep
		err = db.Init(&dml.DBRepository)
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Init DB")
			return nil, nil, nil, err
		}
	}

	writingToFileIsEnabledAndAvailable := flagsGRPC.FileStoragePath != ""
	if flagsGRPC.StoreInterval > 0 && writingToFileIsEnabledAndAvailable {
		if !flagsGRPC.Restore {
			// работаем с новыми метриками, очищая файл со старыми
			err = file.ClearFileContent(flagsGRPC.FileStoragePath)
			if err != nil {
				logger.WriteErrorLog(err.Error(), "ClearFileContent")
			}
		}
		ticker := time.NewTicker(time.Duration(flagsGRPC.StoreInterval) * time.Second)
		go func() {
			err = server.SaveMetricsPerTime(server.MaxSaverWorkTime, ticker, grpcStorage)
			if err != nil {
				logger.WriteErrorLog(err.Error(), "SaveMetricsPerTime")
			}
		}()
	}

	manager := crypto.NewCryptoManager()
	if flagsGRPC.CryptoKey != "" {
		decryptor, err := crypto.NewRSADecryptor(flagsGRPC.CryptoKey)
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Failed to create grpc decryptor")
		}
		manager.SetGRPCDecryptor(decryptor)
	}
	return flagsGRPC, grpcStorage, manager, nil
}

// GetGRPCServer возвращает настроенный и запущенный gRPC сервер
func GetGRPCServer(flags *grpcHandlers.ServerConfigFlags, storage server.Storager, manager *crypto.Manager) (*grpc.Server, error) {
	var err error

	lis, err := net.Listen("tcp", flags.Address)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "TCP listen")
		return nil, err
	}

	trustedIPUnaryInterceptor := interceptors.NewTrustedIPInterceptor(flags.TrustedSubnet)
	decryptUnaryInterceptor := interceptors.NewDecryptUnaryInterceptor(manager)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			trustedIPUnaryInterceptor,
			decryptUnaryInterceptor,
			interceptors.HashCheckUnaryInterceptor,
		),
	)
	pb.RegisterMetricsServer(grpcServer, NewMetricsServer(storage))
	go func() {
		fmt.Println("Server gRPC started")
		if err = grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server Serve error: %v", err)
		}
	}()
	return grpcServer, nil
}
