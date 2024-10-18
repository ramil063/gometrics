package main

import (
	"github.com/ramil063/gometrics/internal/logger"
	"net/http"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/cmd/server/handlers/server"
)

func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}
	var ms = server.NewMemStorage()

	handlers.ParseFlags()
	if err := http.ListenAndServe(handlers.MainURL, server.Router(ms)); err != nil {
		panic(err)
	}
}
