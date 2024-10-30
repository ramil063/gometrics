package handlers

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var MainURL = "localhost:8080"
var StoreInterval = 300
var FileStoragePath = "internal/storage/files/metrics.json"
var Restore = true

type EnvVars struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func ParseFlags() {
	flag.StringVar(&MainURL, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&StoreInterval, "i", 300, "interval of saving metrics to file")
	flag.StringVar(&FileStoragePath, "f", "internal/storage/files/metrics.json", "file storage path")
	flag.BoolVar(&Restore, "r", true, "file storage path")
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

	if ev.Restore != false {
		Restore = ev.Restore
	}
}
