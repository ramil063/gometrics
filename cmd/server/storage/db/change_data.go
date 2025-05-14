package db

import (
	"errors"

	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

// SetGauge создать или обновить метрику типа Gauge
func (s *Storage) SetGauge(name string, value models.Gauge) error {
	result, err := dml.CreateOrUpdateGauge(&dml.DBRepository, name, value)

	if err != nil {
		logger.WriteErrorLog("SetGauge error in sql", err.Error())
		return err
	}
	if result == nil {
		logger.WriteErrorLog("SetGauge error in sql", "empty result")
		return errors.New("SetGauge empty result")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.WriteErrorLog("SetGauge error", err.Error())
		return err
	}
	if rows != 1 {
		logger.WriteErrorLog("SetGauge error", "expected to affect 1 row")
		return errors.New("SetGauge expected to affect 1 row")
	}
	return nil
}

// AddCounter создать или обновить метрику типа Counter
func (s *Storage) AddCounter(name string, value models.Counter) error {
	result, err := dml.CreateOrUpdateCounter(&dml.DBRepository, name, value)

	if err != nil {
		logger.WriteErrorLog("AddCounter error in sql", err.Error())
		return err
	}
	if result == nil {
		logger.WriteErrorLog("AddCounter error in sql", "empty result")
		return errors.New("AddCounter empty result")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.WriteErrorLog("AddCounter error", err.Error())
		return err
	}
	if rows != 1 {
		logger.WriteErrorLog("AddCounter error", "expected to affect 1 row")
		return errors.New("AddCounter expected to affect 1 row")
	}
	return nil
}
