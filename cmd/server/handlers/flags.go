package handlers

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var MainURL = "localhost:8080"

type EnvVars struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func ParseFlags() {
	flag.StringVar(&MainURL, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	var ev EnvVars
	_ = env.Parse(&ev)

	if ev.Address != "" {
		MainURL = ev.Address
	}
}
