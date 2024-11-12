package handlers

import (
	"flag"
	"strconv"

	"github.com/caarlos0/env/v6"

	"github.com/ramil063/gometrics/internal/logger"
)

var MainURL = "localhost:8080"
var StoreInterval = 300
var FileStoragePath = "internal/storage/files/metrics.json"
var Restore = true
var DatabaseDSN = "host=localhost user=u password=pass dbname=db sslmode=disabl"

type EnvVars struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func ParseFlags() {
	flag.StringVar(&MainURL, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&DatabaseDSN, "d", "host=localhost user=u password=pass dbname=db sslmode=disabl", "database DSN")
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

	if !ev.Restore {
		Restore = ev.Restore
	}
	if ev.DatabaseDSN != "" {
		DatabaseDSN = ev.DatabaseDSN
	}

	logger.WriteInfoLog("set g.var", "Address:"+MainURL)
	logger.WriteInfoLog("set g.var", "StoreInterval:"+strconv.Itoa(StoreInterval))
	logger.WriteInfoLog("set g.var", "FileStoragePath:"+FileStoragePath)
	logger.WriteInfoLog("set g.var", "Restore:"+strconv.FormatBool(Restore))
	logger.WriteInfoLog("set g.var", "DatabaseDSN:"+DatabaseDSN)
}
