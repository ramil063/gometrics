package storage

type Gauge float64
type Counter int64

type MemStorageInterface interface {
	SetGauge(name string, value Gauge)
	AddCounter(name string, value Counter)
}

type MemStorage struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
}

func NewMemStorage() MemStorageInterface {
	return &MemStorage{
		Gauges:   make(map[string]Gauge),
		Counters: make(map[string]Counter),
	}
}

func (ms *MemStorage) SetGauge(name string, value Gauge) {
	ms.Gauges[name] = value
}

func (ms *MemStorage) AddCounter(name string, value Counter) {
	oldValue, ok := ms.Counters[name]
	if !ok {
		oldValue = 0
	}
	ms.Counters[name] = oldValue + value
}
