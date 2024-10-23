package main

import (
	"github.com/ramil063/gometrics/cmd/agent/handlers"
	"log"
)

func main() {
	handlers.ParseFlags()
	c := handlers.NewJsonClient()
	r := handlers.NewRequest()
	var err = r.SendMetricsJson(c, 100)
	if err != nil {
		log.Fatal(err)
	}
}
