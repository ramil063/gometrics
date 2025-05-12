package memory

import (
	"errors"
	"sync"

	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

// MemStorage Хранилище данных
type MemStorage struct {
	Gauges   map[string]models.Gauge
	Counters map[string]models.Counter
	mx       sync.RWMutex
}

// StoreGaugeValue сохранение значения метрики типа Gauge
func (ms *MemStorage) StoreGaugeValue(key string, value models.Gauge) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	ms.Gauges[key] = value
}

// GetGaugeValue получение значения метрики типа Gauge по ключу
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

// GetAllGauges получение значений всех метрик типа Gauge
func (ms *MemStorage) GetAllGauges() map[string]models.Gauge {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	mapCopy := make(map[string]models.Gauge, len(ms.Gauges))
	for key, val := range ms.Gauges {
		mapCopy[key] = val
	}
	return mapCopy
}

// StoreCounterValue сохранение значения метрики типа Counter
func (ms *MemStorage) StoreCounterValue(key string, value models.Counter) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	ms.Counters[key] = value
}

// GetCounterValue получение значения метрики типа Counter
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

// GetAllCounters получения значений всех метрик типа Counter
func (ms *MemStorage) GetAllCounters() map[string]models.Counter {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	mapCopy := make(map[string]models.Counter, len(ms.Counters))
	for key, val := range ms.Counters {
		mapCopy[key] = val
	}
	return mapCopy
}

// SetGauge установка значения метрики типа Gauge
func (ms *MemStorage) SetGauge(name string, value models.Gauge) error {
	ms.StoreGaugeValue(name, value)
	return nil
}

// GetGauge получение значения метрики типа Gauge по имени
func (ms *MemStorage) GetGauge(name string) (float64, error) {
	val, err := ms.GetGaugeValue(name)
	return float64(val), err
}

// GetGauges получение значений всех метрик типа Gauge
func (ms *MemStorage) GetGauges() (map[string]models.Gauge, error) {
	return ms.GetAllGauges(), nil
}

// AddCounter добавление(сохранение/обновление) значения метрики типа Counter
func (ms *MemStorage) AddCounter(name string, value models.Counter) error {
	oldValue, err := ms.GetCounterValue(name)
	if err != nil {
		oldValue = 0
		logger.WriteInfoLog("Can't find counter", "AddCounter")
	}
	ms.StoreCounterValue(name, oldValue+value)
	return nil
}

// GetCounter получение значения метрики типа Counter по имени
func (ms *MemStorage) GetCounter(name string) (int64, error) {
	val, err := ms.GetCounterValue(name)
	return int64(val), err
}

// GetCounters получение значений всех метрик типа Counter
func (ms *MemStorage) GetCounters() (map[string]models.Counter, error) {
	return ms.GetAllCounters(), nil
}
