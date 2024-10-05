package storage

type Gauge float64
type Counter int64

type MemStorageInterface interface {
	SetGauge(name string, value Gauge)
	AddCounter(name string, value Counter)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetGauges() map[string]Gauge
	GetCounters() map[string]Counter
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

func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	val, err := ms.Gauges[name]
	return float64(val), err
}
func (ms *MemStorage) GetGauges() map[string]Gauge {
	return ms.Gauges
}

func (ms *MemStorage) AddCounter(name string, value Counter) {
	oldValue, ok := ms.Counters[name]
	if !ok {
		oldValue = 0
	}
	ms.Counters[name] = oldValue + value
}

func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	val, err := ms.Counters[name]
	return int64(val), err
}
func (ms *MemStorage) GetCounters() map[string]Counter {
	return ms.Counters
}
