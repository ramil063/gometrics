package handlers

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/ramil063/gometrics/cmd/agent/config"
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

// CryptoKey путь до публичного ключа шифрования
var CryptoKey = ""

// EnvVars содержит переменные флагов
type EnvVars struct {
	Address        string `env:"ADDRESS"`
	HashKey        string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

// InitFlags парсит глобальные переменные системы, или парсит флаги, или подменяет их значениями по умолчанию
func InitFlags(config *config.AgentConfig) {
	flag.StringVar(&MainURL, "a", config.GetAddress(MainURL), "address and port to run server")
	flag.IntVar(&ReportInterval, "r", config.GetReportInterval(ReportInterval), "report interval in seconds")
	flag.IntVar(&PollInterval, "p", config.GetPollInterval(PollInterval), "poll interval in seconds")
	flag.StringVar(&HashKey, "k", config.GetHashKey(HashKey), "key for hash")
	flag.IntVar(&RateLimit, "l", config.GetRateLimit(RateLimit), "limit requests")
	flag.StringVar(&CryptoKey, "crypto-key", config.GetCryptoKey(CryptoKey), "key for encryption")
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
	if ev.CryptoKey != "" {
		CryptoKey = ev.CryptoKey
	}
}
