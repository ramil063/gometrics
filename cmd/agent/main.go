package main

import (
	"log"

	"github.com/ramil063/gometrics/cmd/agent/handlers"
)

func main() {
	handlers.ParseFlags()
	c := handlers.NewJSONClient()
	r := handlers.NewRequest()
	var err = r.SendMultipleMetricsJSON(c, -1)
	if err != nil {
		log.Fatal(err)
	}
}
