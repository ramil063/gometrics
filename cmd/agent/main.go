package main

import (
	"github.com/ramil063/gometrics/cmd/agent/handlers"
	"log"
)

func main() {
	c := handlers.NewClient()
	r := handlers.NewRequest()
	var err = r.SendMetrics(c, 100)
	if err != nil {
		log.Fatal(err)
	}
}
