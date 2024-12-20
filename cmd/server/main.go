package main

import (
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/logger"
)

func main() {
	var err error
	if err = logger.Initialize(); err != nil {
		panic(err)
	}
	handlers.ParseFlags()

	logger.WriteInfoLog("--------------START SERVER-------------", "")

	var s = server.GetStorage(handlers.Restore, handlers.DatabaseDSN)

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
	if handlers.Restore {
		ticker := time.NewTicker(time.Duration(handlers.StoreInterval) * time.Second)
		go func() {
			if handlers.StoreInterval == 0 {
				return
			}
			server.SaveMetricsPerTime(server.MaxSaverWorkTime, ticker, s)
		}()
	}

	if err = http.ListenAndServe(handlers.MainURL, server.Router(s)); err != nil {
		panic(err)
	}
}
