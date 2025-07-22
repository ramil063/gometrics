package handlers

import (
	"flag"

	"github.com/caarlos0/env/v6"
	serverConfig "github.com/ramil063/gometrics/cmd/server/config"
)

// MainURL основной урл на которым поднят сервис
var MainURL = "localhost:8080"

// StoreInterval на какой промежуток сохраняются данные
var StoreInterval = 300

// FileStoragePath путь для сохранения данных
var FileStoragePath = "internal/storage/files/metrics.json"

// Restore флаг восстановления данных с сохраненного файла
var Restore = true

// DatabaseDSN настройки подключения к БД
var DatabaseDSN = ""

// HashKey ключ для декодирования зашифрованных данных
var HashKey = ""

// CryptoKey путь до приватного ключа шифрования
var CryptoKey = ""

// TrustedSubnet доверенная подсеть для пропуска на сервер
var TrustedSubnet = ""

// EnvVars содержит переменные флагов
type EnvVars struct {
	Address         string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	HashKey         string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
}

// InitFlags парсит глобальные переменные системы, или парсит флаги, или подменяет их значениями по умолчанию
func InitFlags(config *serverConfig.ServerConfig) {
	flag.StringVar(&MainURL, "a", config.GetAddress("localhost:8080"), "address and port to run server")
	flag.StringVar(&DatabaseDSN, "d", config.GetDatabaseDSN(""), "database DSN")
	flag.IntVar(&StoreInterval, "i", config.GetStoreInterval(300), "interval of saving metrics to file")
	flag.StringVar(&FileStoragePath, "f", config.GetFileStoragePath("internal/storage/files/metrics.json"), "file storage path")
	flag.BoolVar(&Restore, "r", config.GetRestore(true), "file storage path")
	flag.StringVar(&HashKey, "k", config.GetHashKey(""), "key for hash")
	flag.StringVar(&CryptoKey, "crypto-key", config.GetCryptoKey(""), "key for encryption")
	flag.StringVar(&TrustedSubnet, "t", config.GetTrustedSubnet(""), "allowed subnet")
	flag.Parse()

	var ev EnvVars
	_ = env.Parse(&ev)

	if ev.Address != "" {
		MainURL = ev.Address
	}

	if ev.StoreInterval != 0 {
		StoreInterval = ev.StoreInterval
	}

	if ev.FileStoragePath != "" {
		FileStoragePath = ev.FileStoragePath
	}

	if !ev.Restore {
		Restore = ev.Restore
	}

	if ev.DatabaseDSN != "" {
		DatabaseDSN = ev.DatabaseDSN
	}

	if ev.HashKey != "" {
		HashKey = ev.HashKey
	}

	if ev.CryptoKey != "" {
		CryptoKey = ev.CryptoKey
	}

	if ev.TrustedSubnet != "" {
		TrustedSubnet = ev.TrustedSubnet
	}

	//only for autotests
	//logger.WriteInfoLog("set g.var", "Address:"+MainURL)
	//logger.WriteInfoLog("set g.var", "StoreInterval:"+strconv.Itoa(StoreInterval))
	//logger.WriteInfoLog("set g.var", "FileStoragePath:"+FileStoragePath)
	//logger.WriteInfoLog("set g.var", "Restore:"+strconv.FormatBool(Restore))
	//logger.WriteInfoLog("set g.var", "DatabaseDSN:"+DatabaseDSN)
	//logger.WriteInfoLog("set g.var", "HashKey:"+HashKey)
}
