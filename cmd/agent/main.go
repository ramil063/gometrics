package main

import (
	"github.com/ramil063/gometrics/cmd/agent/handlers"
)

func main() {
	handlers.ParseFlags()
	c := handlers.NewJSONClient()
	r := handlers.NewRequest()
	r.SendMultipleMetricsJSON(c, -1)
}
