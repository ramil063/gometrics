package main

import (
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/internal/logger"
)

func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}
	handlers.ParseFlags()

	logger.WriteInfoLog("--------------START SERVER-------------", "")

	if err := db.Database.Init(handlers.DatabaseDSN); err != nil {
		panic(err)
	}
	defer db.Database.Ptr.Close()

	var s = server.GetStorage(handlers.Restore)
	ticker := time.NewTicker(time.Duration(handlers.StoreInterval) * time.Second)
	go func() {
		if handlers.StoreInterval == 0 {
			return
		}
		server.SaveMetricsPerTime(server.MaxSaverWorkTime, ticker, s)
	}()

	if err := http.ListenAndServe(handlers.MainURL, server.Router(s)); err != nil {
		panic(err)
	}
}
