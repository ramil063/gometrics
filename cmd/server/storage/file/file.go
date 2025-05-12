package file

import (
	"errors"
	"sync"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

// FStorage хранилище данных
type FStorage struct {
	Gauges   map[string]models.Gauge
	Counters map[string]models.Counter
	mx       sync.RWMutex
}

// StoreGaugeValue сохранение значения метрики типа Gauge
func (s *FStorage) StoreGaugeValue(key string, value models.Gauge) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.Gauges[key] = value
}

// GetGaugeValue получение метрики типа Gauge по ключу
func (s *FStorage) GetGaugeValue(key string) (models.Gauge, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	var err error
	val, ok := s.Gauges[key]
	if !ok {
		err = errors.New("can't get gauge from unknown metric")
	}
	return val, err
}

// GetAllGauges получение значений всех метрик типа Gauge
func (s *FStorage) GetAllGauges() map[string]models.Gauge {
	s.mx.RLock()
	defer s.mx.RUnlock()

	mapCopy := make(map[string]models.Gauge, len(s.Gauges))
	for key, val := range s.Gauges {
		mapCopy[key] = val
	}
	return mapCopy
}

// StoreCounterValue сохранение значения метрики типа Counter
func (s *FStorage) StoreCounterValue(key string, value models.Counter) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.Counters[key] = value
}

// GetCounterValue получение значения метрики типа Counter
func (s *FStorage) GetCounterValue(key string) (models.Counter, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	var err error
	val, ok := s.Counters[key]
	if !ok {
		err = errors.New("can't get counter from unknown metric")
	}
	return val, err
}

// GetAllCounters получение значений всех метрик типа Counter
func (s *FStorage) GetAllCounters() map[string]models.Counter {
	s.mx.RLock()
	defer s.mx.RUnlock()

	mapCopy := make(map[string]models.Counter, len(s.Counters))
	for key, val := range s.Counters {
		mapCopy[key] = val
	}
	return mapCopy
}

// SetGauge установка значения метрики типа Gauge с сохранением в файле
func (s *FStorage) SetGauge(name string, value models.Gauge) error {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "ReadMetricsFromFile SetGauge")
	}
	if metrics == nil {
		metrics = s
	}
	metrics.StoreGaugeValue(name, value)

	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "WriteMetricsToFile SetGauge")
	}
	return err
}

// GetGauge получение значения метрики типа Gauge из файла
func (s *FStorage) GetGauge(name string) (float64, error) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "ReadMetricsFromFile GetGauge")
	}
	if metrics == nil {
		return 0.0, errors.New("no metrics found")
	}
	val, err := metrics.GetGaugeValue(name)

	return float64(val), err
}

// GetGauges получение значений всех метрик типа Gauge из файла
func (s *FStorage) GetGauges() (map[string]models.Gauge, error) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "ReadMetricsFromFile GetGauges")
	}
	if metrics != nil {
		return metrics.GetAllGauges(), nil
	}
	metrics = s
	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "WriteMetricsToFile GetGauges")
		return nil, err
	}
	return metrics.GetAllGauges(), err
}

// AddCounter добавление(сохранение/обновление) значения метрики типа Counter
func (s *FStorage) AddCounter(name string, value models.Counter) error {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "ReadMetricsFromFile AddCounter")
	}
	if metrics == nil {
		metrics = s
		metrics.StoreCounterValue(name, models.Counter(0))
	}

	oldValue, err := metrics.GetCounterValue(name)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "GetCounterValue AddCounter")
	}
	metrics.StoreCounterValue(name, oldValue+value)

	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "WriteMetricsToFile AddCounter")
	}
	return err
}

// GetCounter получение значения метрики типа Counter по имени
func (s *FStorage) GetCounter(name string) (int64, error) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "ReadMetricsFromFile GetCounter")
		return 0, err
	}
	if metrics == nil {
		return 0, errors.New("metrics not found")
	}

	val, err := metrics.GetCounterValue(name)
	if err != nil {
		err = errors.New("can't get counter")
	}
	return int64(val), err
}

// GetCounters получение значений всех метрик типа Counter
func (s *FStorage) GetCounters() (map[string]models.Counter, error) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "ReadMetricsFromFile GetCounters")
	}
	if metrics != nil {
		return metrics.GetAllCounters(), nil
	}

	metrics = s
	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "WriteMetricsToFile GetCounters")
	}
	return metrics.GetAllCounters(), err
}
