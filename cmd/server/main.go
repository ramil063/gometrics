package main

import (
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8080", handlers.Router()); err != nil {
		panic(err)
	}
}
