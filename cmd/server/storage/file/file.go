package file

import (
	"errors"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

type FStorage struct {
	Gauges   map[string]models.Gauge
	Counters map[string]models.Counter
}

func (s *FStorage) SetGauge(name string, value models.Gauge) error {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error read metrics from file", err.Error())
		return err
	}
	if metrics == nil {
		metrics = s
	}
	metrics.Gauges[name] = value

	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error write metrics to file", err.Error())
	}
	return err
}

func (s *FStorage) GetGauge(name string) (float64, error) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error read metrics from file", err.Error())
	}
	if metrics == nil {
		return 0.0, errors.New("no metrics found")
	}
	val, ok := metrics.Gauges[name]
	if !ok {
		err = errors.New("can't set gauge for unknown metric")
	}
	return float64(val), err
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
	if metrics == nil {
		metrics = s
		metrics.Counters[name] = models.Counter(0)
	}

	metrics.Counters[name] += value

	err = WriteMetricsToFile(metrics, handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error write metrics to file", err.Error())
	}
}

func (s *FStorage) GetCounter(name string) (int64, error) {
	metrics, err := ReadMetricsFromFile(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error read metrics from file", err.Error())
	}
	if metrics == nil {
		return 0, errors.New("metrics not found")
	}

	val, ok := metrics.Counters[name]
	if !ok {
		err = errors.New("can't get counter")
	}
	return int64(val), err
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
