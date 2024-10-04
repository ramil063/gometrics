package main

import (
	"github.com/ramil063/gometrics/cmd/server/handlers"
)

func main() {
	server := handlers.NewServer(":8080")
	if err := server.Run(); err != nil {
		panic(err)
	}
}
