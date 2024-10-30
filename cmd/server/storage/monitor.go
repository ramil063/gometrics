package storage

import (
	"github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/cmd/server/filer"
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/logger"
	"go.uber.org/zap"
	"log"
	"time"
)

var workSecond = 0

func GetMonitor(restore bool) storage.Monitor {
	var m = &storage.Monitor{}
	if restore {
		Reader, err := filer.NewReader(handlers.FileStoragePath)
		if err != nil {
			logger.Log.Error("error create monitor reader", zap.Error(err))
		}
		defer Reader.Close()

		m, err = Reader.ReadMonitor()
		if err != nil {
			temp := storage.NewMonitor()
			temp.PollCount = 0
			m = &temp
			logger.Log.Error("error read monitor", zap.Error(err))
		}
		return *m
	}
	res := storage.NewMonitor()
	res.PollCount = 0
	return res
}

// SaveMonitorPerSeconds сохранение метрик в единицу времени
func SaveMonitorPerSeconds(workTime int, ticker *time.Ticker, storeInterval int, filePath string) error {

	for workSecond < workTime {
		<-ticker.C
		workSecond++
		if (workSecond % storeInterval) == 0 {
			log.Println("save metrics")
			m := storage.NewMonitor()
			Writer, err := filer.NewWriter(filePath)
			if err != nil {
				logger.Log.Error("error create monitor writer", zap.Error(err))
				return err
			}
			defer Writer.Close()

			err = Writer.WriteMonitor(&m)
			if err != nil {
				logger.Log.Error("error write monitor", zap.Error(err))
				return err
			}
		}
	}
	return nil
}
