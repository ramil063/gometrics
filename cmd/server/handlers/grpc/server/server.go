package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	grpcServerConfig "github.com/ramil063/gometrics/cmd/server/config/grpc"
	grpcHandlers "github.com/ramil063/gometrics/cmd/server/handlers/grpc"
	"github.com/ramil063/gometrics/cmd/server/handlers/grpc/interceptors"
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/cmd/server/storage/file"
	pb "github.com/ramil063/gometrics/internal/grpc/proto"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

// GetGRPCServer возвращает настроенный и запущенный gRPC сервер
func GetGRPCServer() (*grpc.Server, error) {
	configGRPC, err := grpcServerConfig.GetConfig()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "GetConfig")
		return nil, err
	}
	flagsGRPC, err := grpcHandlers.GetFlags(configGRPC)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "GetFlags")
		return nil, err
	}

	var grpcStorage = server.GetStorage(flagsGRPC.FileStoragePath, flagsGRPC.DatabaseDSN)

	if flagsGRPC.DatabaseDSN != "" {
		rep, errRepo := dml.NewRepository()
		if errRepo != nil {
			logger.WriteErrorLog(errRepo.Error(), "NewRepository")
			return nil, errRepo
		}
		dml.DBRepository = *rep
		err = db.Init(&dml.DBRepository)
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Init DB")
			return nil, err
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

	if flagsGRPC.CryptoKey != "" {
		crypto.GRPCDecryptor, err = crypto.NewRSADecryptor(flagsGRPC.CryptoKey)
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Failed to create grpc decryptor")
		}
	}

	lis, err := net.Listen("tcp", flagsGRPC.Address)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "TCP listen")
		return nil, err
	}

	trustedIPUnaryInterceptor := interceptors.NewTrustedIPInterceptor(flagsGRPC.TrustedSubnet)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			trustedIPUnaryInterceptor,
			interceptors.DecryptUnaryInterceptor,
			interceptors.HashCheckUnaryInterceptor,
		),
	)
	pb.RegisterMetricsServer(grpcServer, NewMetricsServer(grpcStorage))
	go func() {
		fmt.Println("Server gRPC started")
		if err = grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server Serve error: %v", err)
		}
	}()
	return grpcServer, nil
}
