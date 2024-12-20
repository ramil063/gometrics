package memory

import (
	"errors"

	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

type MemStorage struct {
	Gauges   map[string]models.Gauge
	Counters map[string]models.Counter
}

func (ms *MemStorage) SetGauge(name string, value models.Gauge) error {
	ms.Gauges[name] = value
	return nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	val, ok := ms.Gauges[name]
	var err error
	if !ok {
		err = errors.New("gauge not found")
	}
	return float64(val), err
}

func (ms *MemStorage) GetGauges() (map[string]models.Gauge, error) {
	return ms.Gauges, nil
}

func (ms *MemStorage) AddCounter(name string, value models.Counter) error {
	oldValue, ok := ms.Counters[name]
	if !ok {
		oldValue = 0
		logger.WriteInfoLog("Can't find counter", "AddCounter")
	}
	ms.Counters[name] = oldValue + value
	return nil
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
	val, ok := ms.Counters[name]
	var err error
	if !ok {
		err = errors.New("counter not found")
	}
	return int64(val), err
}

func (ms *MemStorage) GetCounters() (map[string]models.Counter, error) {
	return ms.Counters, nil
}
