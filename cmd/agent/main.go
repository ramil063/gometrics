package main

import (
	"github.com/ramil063/gometrics/cmd/agent/handlers"
	"log"
)

func main() {
	handlers.ParseFlags()
	c := handlers.NewJSONClient()
	r := handlers.NewRequest()
	var err = r.SendMetricsJSON(c, 100)
	if err != nil {
		log.Fatal(err)
	}
}
