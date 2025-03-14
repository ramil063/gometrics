package server

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

var workSecond = 0

// SaveMetricsPerTime сохранение метрик в единицу времени
func SaveMetricsPerTime(workTime int, ticker *time.Ticker, s Storager) error {
	quit := make(chan struct{})
	for workSecond < workTime {
		select {
		case <-ticker.C:
			workSecond++
			m := storage.NewMonitor()
			err := PrepareMetricsValues(s, m)
			if err != nil {
				ticker.Stop()
				return err
			}
		case <-quit:
			ticker.Stop()
		}
	}
	return nil
}

func PrepareMetricsValues(s Storager, m *storage.Monitor) error {
	v := reflect.ValueOf(m).Elem()
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		metricID := typeOfS.Field(i).Name
		if metricID == "mx" {
			continue
		}
		metricValue, _ := strconv.ParseFloat(fmt.Sprintf("%v", v.Field(i).Interface()), 64)

		if typeOfS.Field(i).Name == "PollCount" {
			err := s.AddCounter(metricID, models.Counter(1))
			if err != nil {
				logger.WriteErrorLog(err.Error(), "Counter")
				return err
			}
		} else {
			err := s.SetGauge(metricID, models.Gauge(metricValue))
			if err != nil {
				logger.WriteErrorLog(err.Error(), "Gauge")
				return err
			}
		}
	}
	return nil
}
