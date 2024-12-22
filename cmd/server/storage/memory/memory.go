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

func (s *MemStorage) StoreGaugeValue(key string, value models.Gauge) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.Gauges[key] = value
}

func (s *MemStorage) GetGaugeValue(key string) (models.Gauge, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	var err error
	val, ok := s.Gauges[key]
	if !ok {
		err = errors.New("can't set gauge for unknown metric")
	}
	return val, err
}

func (s *MemStorage) GetAllGauges() map[string]models.Gauge {
	s.mx.RLock()
	defer s.mx.RUnlock()

	mapCopy := make(map[string]models.Gauge, len(s.Gauges))
	for key, val := range s.Gauges {
		mapCopy[key] = val
	}
	return mapCopy
}

func (s *MemStorage) StoreCounterValue(key string, value models.Counter) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.Counters[key] = value
}

func (s *MemStorage) GetCounterValue(key string) (models.Counter, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	var err error
	val, ok := s.Counters[key]
	if !ok {
		err = errors.New("can't set counter for unknown metric")
	}
	return val, err
}

func (s *MemStorage) GetAllCounters() map[string]models.Counter {
	s.mx.RLock()
	defer s.mx.RUnlock()

	mapCopy := make(map[string]models.Counter, len(s.Counters))
	for key, val := range s.Counters {
		mapCopy[key] = val
	}
	return mapCopy
}

func (ms *MemStorage) SetGauge(name string, value models.Gauge) error {
	ms.Gauges[name] = value
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
