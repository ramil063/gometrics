package file

import (
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

type FStorage struct {
	Gauges   map[string]models.Gauge
	Counters map[string]models.Counter
}

func (s *FStorage) SetGauge(name string, value models.Gauge) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error read metrics from file", err.Error())
	}
	if metrics == nil {
		metrics = s
	}
	metrics.Gauges[name] = value

	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error write metrics to file", err.Error())
	}
}

func (s *FStorage) GetGauge(name string) (float64, bool) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error read metrics from file", err.Error())
	}
	if metrics == nil {
		return 0.0, false
	}
	val, ok := metrics.Gauges[name]
	return float64(val), ok
}

func (s *FStorage) GetGauges() map[string]models.Gauge {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error read metrics from file", err.Error())
	}
	if metrics != nil {
		return metrics.Gauges
	}
	metrics = s
	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error write metrics to file", err.Error())
	}
	return metrics.Gauges
}

func (s *FStorage) AddCounter(name string, value models.Counter) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error read metrics from file", err.Error())
	}
	oldValue := models.Counter(0)
	if metrics != nil {
		oldValue = metrics.Counters[name]
	}

	metrics.Counters[name] = oldValue + value

	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error write metrics to file", err.Error())
	}
}

func (s *FStorage) GetCounter(name string) (int64, bool) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error read metrics from file", err.Error())
	}
	if metrics == nil {
		return 0, false
	}

	val, ok := metrics.Counters[name]
	return int64(val), ok
}

func (s *FStorage) GetCounters() map[string]models.Counter {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error read metrics from file", err.Error())
	}
	if metrics != nil {
		return metrics.Counters
	}

	metrics = s
	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error write metrics to file", err.Error())
	}
	return metrics.Counters
}
