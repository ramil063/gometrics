package handlers

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

var MainURL = "localhost:8080"
var PollInterval = 2
var ReportInterval = 10
var HashKey = ""
var RateLimit = 1

type EnvVars struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	HashKey        string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func ParseFlags() {
	flag.StringVar(&MainURL, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&PollInterval, "p", 2, "poll interval in seconds")
	flag.StringVar(&HashKey, "k", "", "key for hash")
	flag.IntVar(&RateLimit, "l", 1, "limit requests")
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
	if ev.HashKey != "" {
		HashKey = ev.HashKey
	}
	if ev.RateLimit != 0 {
		RateLimit = ev.RateLimit
	}
}
