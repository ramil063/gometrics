package handlers

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// MainURL основной урл на который нужно отправлять метрики
var MainURL = "localhost:8080"

// PollInterval с каким интервалом в секундах нужно собирать метрики
var PollInterval = 2

// ReportInterval с каким интервалом в секундах нужно отправлять данные на удаленный сервис
var ReportInterval = 10

// HashKey ключ для шифрования и дешифровки передаваемых данных
var HashKey = ""

// RateLimit количество одновременных запросов отправляемых на удаленный сервис
var RateLimit = 1

// EnvVars содержит переменные флагов
type EnvVars struct {
	Address        string `env:"ADDRESS"`
	HashKey        string `env:"KEY"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

// ParseFlags парсит глобальные переменные системы, или парсит флаги, или подменяет их значениями по умолчанию
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
