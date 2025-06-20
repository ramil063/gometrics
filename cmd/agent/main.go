package main

import (
	"fmt"

	"github.com/ramil063/gometrics/cmd/agent/handlers"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	handlers.ParseFlags()

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	c := handlers.NewJSONClient()
	r := handlers.NewRequest()
	r.SendMultipleMetricsJSON(c, -1)
}
