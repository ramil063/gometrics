package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

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
	handlers.ParseFlags()

	if handlers.CryptoKey != "" {
		var err error
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
		rep, err := dml.NewRepository()
		if err != nil {
			logger.WriteErrorLog(err.Error(), "NewRepository")
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
			err := file.ClearFileContent(handlers.FileStoragePath)
			if err != nil {
				logger.WriteErrorLog(err.Error(), "ClearFileContent")
			}
		}
		ticker := time.NewTicker(time.Duration(handlers.StoreInterval) * time.Second)
		go func() {
			server.SaveMetricsPerTime(server.MaxSaverWorkTime, ticker, s)
		}()
	}
	if err := http.ListenAndServe(handlers.MainURL, server.Router(s)); err != nil {
		panic(err)
	}
}
