package main

import (
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"net/http"
)

func main() {
	handlers.ParseFlags()
	if err := http.ListenAndServe(handlers.MainURL, handlers.Router()); err != nil {
		panic(err)
	}
}
