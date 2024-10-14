package storage

import (
	"github.com/ramil063/gometrics/cmd/server/models"
)

type MemStorage struct {
	Gauges   map[string]models.Gauge
	Counters map[string]models.Counter
}

func (ms *MemStorage) SetGauge(name string, value models.Gauge) {
	ms.Gauges[name] = value
}

func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	val, ok := ms.Gauges[name]
	return float64(val), ok
}
func (ms *MemStorage) GetGauges() map[string]models.Gauge {
	return ms.Gauges
}

func (ms *MemStorage) AddCounter(name string, value models.Counter) {
	oldValue, ok := ms.Counters[name]
	if !ok {
		oldValue = 0
	}
	ms.Counters[name] = oldValue + value
}

func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	val, ok := ms.Counters[name]
	return int64(val), ok
}
func (ms *MemStorage) GetCounters() map[string]models.Counter {
	return ms.Counters
}
