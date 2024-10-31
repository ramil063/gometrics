package storage

import (
	"log"
	"reflect"
	"time"

	"github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/cmd/server/filer"
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

var workSecond = 0

func GetMonitor(restore bool) storage.Monitor {
	var m = &storage.Monitor{}
	if restore {
		Reader, err := filer.NewReader(handlers.FileStoragePath)
		if err != nil {
			logger.WriteErrorLog("error create monitor reader", err.Error())
		}
		defer Reader.Close()

		m, err = Reader.ReadMonitor()
		if err != nil {
			logger.WriteInfoLog("error in read monitor", err.Error())
			if m == nil {
				temp := storage.NewMonitor()
				temp.PollCount = 0
				m = &temp
			}
			Writer, err := filer.NewWriter(handlers.FileStoragePath)
			if err != nil {
				logger.WriteInfoLog("error create file writer", err.Error())
			}
			defer Writer.Close()
			err = Writer.WriteMonitor(m)
			if err != nil {
				logger.WriteInfoLog("error write monitor", err.Error())
			}
		}
		return *m
	}
	res := storage.NewMonitor()
	res.PollCount = 0
	return res
}

// SaveMonitorPerSeconds сохранение метрик в единицу времени
func SaveMonitorPerSeconds(workTime int, ticker *time.Ticker, filePath string) error {
	for workSecond < workTime {
		<-ticker.C
		workSecond++
		return SaveMonitor(filePath)
	}
	return nil
}

// SaveMonitor сохранение метрик
func SaveMonitor(filePath string) error {
	log.Println("save monitor")
	m := storage.NewMonitor()
	Writer, err := filer.NewWriter(filePath)
	if err != nil {
		logger.WriteErrorLog("error create monitor writer", err.Error())
		return err
	}
	defer Writer.Close()

	err = Writer.WriteMonitor(&m)
	if err != nil {
		logger.WriteErrorLog("error write monitor", err.Error())
		return err
	}
	return nil
}

// SaveMetric сохранение метрики
func SaveMetric(name string, mType string, value models.Gauge, delta models.Counter, filePath string) error {
	Reader, err := filer.NewReader(filePath)
	if err != nil {
		logger.WriteErrorLog("error create monitor reader", err.Error())
	}
	defer Reader.Close()

	m, err := Reader.ReadMonitor()
	if err != nil {
		logger.WriteInfoLog("error in read monitor", err.Error())
	}

	if mType == "counter" {
		oldValue := reflect.ValueOf(m).Elem().FieldByName(name)
		reflect.ValueOf(m).Elem().FieldByName(name).SetInt(oldValue.Int() + int64(delta))
	} else {
		oldValue := reflect.ValueOf(m).Elem().FieldByName(name)
		reflect.ValueOf(m).Elem().FieldByName(name).SetFloat(oldValue.Float() + float64(value))
	}

	Writer, err := filer.NewWriter(filePath)
	if err != nil {
		logger.WriteErrorLog("error create monitor writer", err.Error())
		return err
	}
	defer Writer.Close()

	err = Writer.WriteMonitor(m)
	if err != nil {
		logger.WriteErrorLog("error write monitor", err.Error())
		return err
	}
	return nil
}
