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
	var ms = server.NewMemStorage()
	m := storage.GetMonitor(handlers.Restore)
	server.PrepareStorageValues(ms, m)

	ticker := time.NewTicker(time.Second)
	go storage.SaveMonitorPerSeconds(server.MaxSaverWorkTime, ticker, handlers.StoreInterval, handlers.FileStoragePath)

	handlers.ParseFlags()
	if err := http.ListenAndServe(handlers.MainURL, server.Router(ms)); err != nil {
		panic(err)
	}
}
