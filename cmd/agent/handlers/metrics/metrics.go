package metrics

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

// CollectMetricsRequestBodies сбор метрик в тела для отправки на сторонний сервис
func CollectMetricsRequestBodies(monitor *storage.Monitor) []byte {
	allMetrics := GetMetricsCollection(monitor)
	body, err := json.Marshal(allMetrics)
	if err != nil {
		logger.WriteErrorLog("Error marshal metrics", err.Error())
	}
	return body
}

// GetMetricsCollection сбор метрик для отправки на сторонний сервис
func GetMetricsCollection(monitor *storage.Monitor) []models.Metrics {
	v := reflect.ValueOf(monitor).Elem()
	typeOfS := v.Type()
	allMetrics := make([]models.Metrics, 0, 100)

	for i := 0; i < v.NumField(); i++ {
		metricID := typeOfS.Field(i).Name

		if metricID == "mx" {
			continue
		}

		metricValue, _ := strconv.ParseFloat(fmt.Sprintf("%v", v.Field(i).Interface()), 64)
		delta := int64(monitor.GetCountValue())

		if metricID == "CPUutilization" {
			CPUutilization := monitor.GetAllCPUutilization()
			for j, value := range CPUutilization {
				valuePtr := float64(value)
				metrics := models.Metrics{
					ID:    "CPUutilization" + strconv.Itoa(j),
					MType: "gauge",
					Delta: nil,
					Value: &valuePtr,
				}
				allMetrics = append(allMetrics, metrics)
			}
			continue
		}
		metrics := models.Metrics{
			ID:    metricID,
			MType: "gauge",
			Delta: nil,
			Value: &metricValue,
		}

		if typeOfS.Field(i).Name == "PollCount" {
			metrics.MType = "counter"
			metrics.Delta = &delta
			metrics.Value = nil
		}
		allMetrics = append(allMetrics, metrics)
	}
	return allMetrics
}

// CollectMonitorMetrics собирает метрики монитора
func CollectMonitorMetrics(count int, monitor *storage.Monitor, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		storage.SetMetricsToMonitor(monitor)
		monitor.StoreCountValue(count)
	}()
}

// CollectGopsutilMetrics собирает метрики через gopsutil
func CollectGopsutilMetrics(monitor *storage.Monitor, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := storage.SetGopsutilMetricsToMonitor(monitor)
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Error in set gopsutil metrics")
		}
	}()
}
