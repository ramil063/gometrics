package main

import (
	"net/http"

	agentStorage "github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
	"github.com/ramil063/gometrics/internal/logger"
)

func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}
	var ms = server.NewMemStorage()
	m := agentStorage.NewMonitor()
	server.PrepareStorageValues(ms, m)

	handlers.ParseFlags()
	if err := http.ListenAndServe(handlers.MainURL, server.Router(ms)); err != nil {
		panic(err)
	}
}
