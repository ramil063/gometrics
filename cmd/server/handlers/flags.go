package handlers

import (
	"flag"

	"github.com/caarlos0/env/v6"
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

// EnvVars содержит переменные флагов
type EnvVars struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	HashKey         string `env:"KEY"`
}

// ParseFlags парсит глобальные переменные системы, или парсит флаги, или подменяет их значениями по умолчанию
func ParseFlags() {
	flag.StringVar(&MainURL, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&DatabaseDSN, "d", "", "database DSN")
	flag.IntVar(&StoreInterval, "i", 300, "interval of saving metrics to file")
	flag.StringVar(&FileStoragePath, "f", "internal/storage/files/metrics.json", "file storage path")
	flag.BoolVar(&Restore, "r", true, "file storage path")
	flag.StringVar(&HashKey, "k", "", "key for hash")
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

	//only for autotests
	//logger.WriteInfoLog("set g.var", "Address:"+MainURL)
	//logger.WriteInfoLog("set g.var", "StoreInterval:"+strconv.Itoa(StoreInterval))
	//logger.WriteInfoLog("set g.var", "FileStoragePath:"+FileStoragePath)
	//logger.WriteInfoLog("set g.var", "Restore:"+strconv.FormatBool(Restore))
	//logger.WriteInfoLog("set g.var", "DatabaseDSN:"+DatabaseDSN)
	//logger.WriteInfoLog("set g.var", "HashKey:"+HashKey)
}
