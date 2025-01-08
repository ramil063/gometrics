package server

import (
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/cmd/server/storage/file"
	"github.com/ramil063/gometrics/cmd/server/storage/memory"
	"github.com/ramil063/gometrics/internal/models"
)

func NewMemStorage() Storager {
	return &memory.MemStorage{
		Gauges:   make(map[string]models.Gauge),
		Counters: make(map[string]models.Counter),
	}
}

func NewFileStorage() Storager {
	return &file.FStorage{
		Gauges:   make(map[string]models.Gauge),
		Counters: make(map[string]models.Counter),
	}
}

func NewDBStorage() Storager {
	return &db.Storage{}
}

// GetStorage получить хранителя данных
func GetStorage(fileStoragePath string, dsn string) Storager {
	if dsn != "" {
		return NewDBStorage()
	}
	if fileStoragePath != "" {
		return NewFileStorage()
	}
	return NewMemStorage()
}
