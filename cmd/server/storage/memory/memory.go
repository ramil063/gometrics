package memory

import (
	"errors"
	"sync"

	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

type MemStorage struct {
	Gauges   map[string]models.Gauge
	Counters map[string]models.Counter
	mx       sync.RWMutex
}

func (ms *MemStorage) StoreGaugeValue(key string, value models.Gauge) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	ms.Gauges[key] = value
}

func (ms *MemStorage) GetGaugeValue(key string) (models.Gauge, error) {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	var err error
	val, ok := ms.Gauges[key]
	if !ok {
		err = errors.New("can't set gauge for unknown metric")
	}
	return val, err
}

func (ms *MemStorage) GetAllGauges() map[string]models.Gauge {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	mapCopy := make(map[string]models.Gauge, len(ms.Gauges))
	for key, val := range ms.Gauges {
		mapCopy[key] = val
	}
	return mapCopy
}

func (ms *MemStorage) StoreCounterValue(key string, value models.Counter) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	ms.Counters[key] = value
}

func (ms *MemStorage) GetCounterValue(key string) (models.Counter, error) {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	var err error
	val, ok := ms.Counters[key]
	if !ok {
		err = errors.New("can't set counter for unknown metric")
	}
	return val, err
}

func (ms *MemStorage) GetAllCounters() map[string]models.Counter {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	mapCopy := make(map[string]models.Counter, len(ms.Counters))
	for key, val := range ms.Counters {
		mapCopy[key] = val
	}
	return mapCopy
}

func (ms *MemStorage) SetGauge(name string, value models.Gauge) error {
	ms.StoreGaugeValue(name, value)
	return nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	val, err := ms.GetGaugeValue(name)
	return float64(val), err
}

func (ms *MemStorage) GetGauges() (map[string]models.Gauge, error) {
	return ms.GetAllGauges(), nil
}

func (ms *MemStorage) AddCounter(name string, value models.Counter) error {
	oldValue, err := ms.GetCounterValue(name)
	if err != nil {
		oldValue = 0
		logger.WriteInfoLog("Can't find counter", "AddCounter")
	}
	ms.StoreCounterValue(name, oldValue+value)
	return nil
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
	val, err := ms.GetCounterValue(name)
	return int64(val), err
}

func (ms *MemStorage) GetCounters() (map[string]models.Counter, error) {
	return ms.GetAllCounters(), nil
}
