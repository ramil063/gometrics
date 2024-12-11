package db

import (
	"context"

	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

func (s *Storage) GetGauge(name string) (float64, bool) {
	row := dml.DBRepository.QueryRowContext(context.Background(), "SELECT value FROM gauge WHERE name = $1", name)
	var selectedValue float64

	err := row.Scan(&selectedValue)

	return selectedValue, err == nil
}

func (s *Storage) GetGauges() map[string]models.Gauge {
	result := make(map[string]models.Gauge)

	rows, err := dml.DBRepository.QueryContext(context.Background(), "SELECT name, value FROM gauge")
	if err != nil {
		return result
	}
	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	var name string
	var value float64
	// пробегаем по всем записям
	for rows.Next() {
		err = rows.Scan(&name, &value)
		if err != nil {
			logger.WriteErrorLog("GetGauges error in sql", err.Error())
			continue
		}
		result[name] = models.Gauge(value)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		logger.WriteErrorLog("GetGauges error in rows", err.Error())
	}
	return result
}

func (s *Storage) GetCounter(name string) (int64, bool) {
	row := dml.DBRepository.QueryRowContext(context.Background(), "SELECT value FROM counter WHERE name = $1", name)
	var selectedValue int64
	err := row.Scan(&selectedValue)

	return selectedValue, err == nil
}

func (s *Storage) GetCounters() map[string]models.Counter {
	result := make(map[string]models.Counter)
	rows, err := dml.DBRepository.QueryContext(context.Background(), "SELECT name, value FROM counter")
	if err != nil {
		return result
	}
	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	var name string
	var value int64
	// пробегаем по всем записям
	for rows.Next() {
		err = rows.Scan(&name, &value)
		if err != nil {
			logger.WriteErrorLog("GetCounters error in sql", err.Error())
			continue
		}
		result[name] = models.Counter(value)
	}
	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		logger.WriteErrorLog("GetCounters error in rows", err.Error())
	}
	return result
}
