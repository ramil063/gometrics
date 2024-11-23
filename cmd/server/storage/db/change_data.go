package db

import (
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

func (s *Storage) SetGauge(name string, value models.Gauge) {
	result, err := dml.CreateOrUpdateGauge(&dml.DBRepository, name, value)

	if err != nil {
		logger.WriteErrorLog("SetGauge error in sql", err.Error())
		return
	}
	if result == nil {
		logger.WriteErrorLog("SetGauge error in sql", "empty result")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.WriteErrorLog("SetGauge error", err.Error())
	}
	if rows != 1 {
		logger.WriteErrorLog("SetGauge error", "expected to affect 1 row")
	}
}

func (s *Storage) AddCounter(name string, value models.Counter) {
	result, err := dml.CreateOrUpdateCounter(&dml.DBRepository, name, value)

	if err != nil {
		logger.WriteErrorLog("AddCounter error in sql", err.Error())
		return
	}
	if result == nil {
		logger.WriteErrorLog("AddCounter error in sql", "empty result")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.WriteErrorLog("AddCounter error", err.Error())
	}
	if rows != 1 {
		logger.WriteErrorLog("AddCounter error", "expected to affect 1 row")
	}
}
