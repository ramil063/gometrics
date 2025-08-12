package grpc

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/ramil063/gometrics/cmd/agent/config"
)

// SystemConfigFlags содержит переменные флагов
// Address основной урл на который нужно отправлять метрики
// PollInterval с каким интервалом в секундах нужно собирать метрики
// ReportInterval с каким интервалом в секундах нужно отправлять данные на удаленный сервис
// HashKey ключ для шифрования и дешифровки передаваемых данных
// RateLimit количество одновременных запросов отправляемых на удаленный сервис
// CryptoKey путь до публичного ключа шифрования
type SystemConfigFlags struct {
	Address        string `env:"GRPC_ADDRESS"`
	HashKey        string `env:"GRPC_KEY"`
	CryptoKey      string `env:"GRPC_CRYPTO_KEY"`
	ReportInterval int    `env:"GRPC_REPORT_INTERVAL"`
	PollInterval   int    `env:"GRPC_POLL_INTERVAL"`
	RateLimit      int    `env:"GRPC_RATE_LIMIT"`
}

// GetFlags парсит глобальные переменные системы, или парсит флаги, или подменяет их значениями по умолчанию
func GetFlags(config *config.AgentConfig) (*SystemConfigFlags, error) {

	//значения флагов по умолчанию
	flags := &SystemConfigFlags{
		Address:        "localhost:3202",
		PollInterval:   2,
		ReportInterval: 10,
		RateLimit:      1,
	}

	var (
		address        string
		hashKey        string
		cryptoKey      string
		reportInterval int
		pollInterval   int
		rateLimit      int
	)

	flag.StringVar(&address, "grpc-a", config.GetAddress(flags.Address), "address and port to run server")
	flag.IntVar(&reportInterval, "grpc-r", config.GetReportInterval(flags.ReportInterval), "report interval in seconds")
	flag.IntVar(&pollInterval, "grpc-p", config.GetPollInterval(flags.PollInterval), "poll interval in seconds")
	flag.StringVar(&hashKey, "grpc-k", config.GetHashKey(flags.HashKey), "key for hash")
	flag.IntVar(&rateLimit, "grpc-l", config.GetRateLimit(flags.RateLimit), "limit requests")
	flag.StringVar(&cryptoKey, "grpc-crypto-key", config.GetCryptoKey(flags.CryptoKey), "key for encryption")
	flag.Parse()

	var envVars SystemConfigFlags
	err := env.Parse(&envVars)
	if err != nil {
		return flags, fmt.Errorf("error parsing environment variables: %w", err)
	}

	applyFlags(flags, address, reportInterval, pollInterval, hashKey, rateLimit, cryptoKey)
	applyEnvVars(flags, envVars)

	return flags, nil
}

// applyFlags присваивание флагов переданных в командной строке
func applyFlags(flags *SystemConfigFlags, address string, reportInterval, pollInterval int, hashKey string, rateLimit int, cryptoKey string) {
	if address != "" && address != flags.Address {
		flags.Address = address
	}
	if reportInterval != 0 && reportInterval != flags.ReportInterval {
		flags.ReportInterval = reportInterval
	}
	if pollInterval != 0 && pollInterval != flags.PollInterval {
		flags.PollInterval = pollInterval
	}
	if hashKey != "" && hashKey != flags.HashKey {
		flags.HashKey = hashKey
	}
	if rateLimit != 0 && rateLimit != flags.RateLimit {
		flags.RateLimit = rateLimit
	}
	if cryptoKey != "" && cryptoKey != flags.CryptoKey {
		flags.CryptoKey = cryptoKey
	}
}

// applyEnvVars присваивание переменных окружения
func applyEnvVars(flags *SystemConfigFlags, envVars SystemConfigFlags) {
	if envVars.Address != "" {
		flags.Address = envVars.Address
	}
	if envVars.ReportInterval != 0 {
		flags.ReportInterval = envVars.ReportInterval
	}
	if envVars.PollInterval != 0 {
		flags.PollInterval = envVars.PollInterval
	}
	if envVars.HashKey != "" {
		flags.HashKey = envVars.HashKey
	}
	if envVars.RateLimit != 0 {
		flags.RateLimit = envVars.RateLimit
	}
	if envVars.CryptoKey != "" {
		flags.CryptoKey = envVars.CryptoKey
	}
}
