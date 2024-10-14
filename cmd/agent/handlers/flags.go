package handlers

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

var MainURL = "localhost:8080"
var PollInterval = 2
var ReportInterval = 10

type EnvVars struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func ParseFlags() {
	flag.StringVar(&MainURL, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&PollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()

	var ev EnvVars
	_ = env.Parse(&ev)

	if ev.Address != "" {
		MainURL = ev.Address
	}
	if ev.ReportInterval != 0 {
		ReportInterval = ev.ReportInterval
	}
	if ev.PollInterval != 0 {
		PollInterval = ev.PollInterval
	}
}
