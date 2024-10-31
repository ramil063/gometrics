package main

import (
	"net/http"
	"time"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	"github.com/ramil063/gometrics/cmd/server/storage"
	"github.com/ramil063/gometrics/internal/logger"
)

func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}

	logger.WriteInfoLog("--------------START SERVER-------------", "")

	var ms = server.NewMemStorage()
	m := storage.GetMonitor(handlers.Restore)
	server.PrepareStorageValues(ms, m)

	ticker := time.NewTicker(time.Duration(handlers.StoreInterval) * time.Second)
	go func() {
		if handlers.StoreInterval == 0 {
			return
		}
		err := storage.SaveMonitorPerSeconds(server.MaxSaverWorkTime, ticker, handlers.FileStoragePath)
		if err != nil {
			logger.WriteInfoLog("error in SaveMonitorPerSeconds", err.Error())
		}
	}()

	handlers.ParseFlags()
	if err := http.ListenAndServe(handlers.MainURL, server.Router(ms)); err != nil {
		panic(err)
	}
}
