package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	metricsHandler "github.com/ramil063/gometrics/cmd/agent/handlers/metrics"
	"github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/internal/errors"
	pb "github.com/ramil063/gometrics/internal/grpc/proto"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

// Clienter работа с клиентом в формате json
type Clienter interface {
	Close() error
	SendMetrics(ctx context.Context, metrics []*pb.Metric, encryptedMetrics []byte) error
}

// Requester отправляет данные
type Requester interface {
	SendMetricsProcess(c Clienter, maxCount int, ctxGrSh context.Context, flags *SystemConfigFlags, manager *crypto.Manager)
}

type request struct {
	IP string
}

func NewRequest() Requester {
	req := request{}
	ip, err := req.getOutboundIP()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "IP")
	}
	req.IP = ip
	return req
}

func (r request) getOutboundIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	addrs, err := net.LookupHost(hostname)
	if err != nil || len(addrs) == 0 {
		return "", fmt.Errorf("cannot determine IP")
	}
	return addrs[0], nil
}

// SendMetricsProcess отправка нескольких метрик
func (r request) SendMetricsProcess(c Clienter, maxCount int, ctxGrSh context.Context, flags *SystemConfigFlags, manager *crypto.Manager) {
	var pollInterval = time.Duration(flags.PollInterval) * time.Second
	var reportInterval = time.Duration(flags.ReportInterval) * time.Second
	count := 0

	var monitor storage.Monitor
	tickerPool := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)

	var sendMonitor = make(chan *storage.Monitor, 1)
	defer close(sendMonitor)

	var mu sync.Mutex
	var shutdown = false

	log.Println("grpc agent start")

	go func() {
		defer tickerPool.Stop()
		for maxCount < 0 {
			<-tickerPool.C

			mu.Lock()
			count++
			mu.Unlock()

			log.Println("get metrics grpc start")
			var collectWg sync.WaitGroup
			metricsHandler.CollectMonitorMetrics(count, &monitor, &collectWg)
			metricsHandler.CollectGopsutilMetrics(&monitor, &collectWg)
			collectWg.Wait()

			sendMonitor <- &monitor
			log.Println("get metrics grpc end")
			if len(sendMonitor) == 1 {
				<-sendMonitor
			}
		}
	}()

	go func() {
		defer tickerReport.Stop()
		for maxCount < 0 {
			<-tickerReport.C

			log.Println("send metrics grpc start")
			mon := <-sendMonitor
			log.Println("send metrics grpc count value=", mon.GetCountValue())

			for worker := 0; worker < flags.RateLimit; worker++ {
				log.Println("send metrics grpc worker=", worker)
				go SendMetricsByGRPC(r, c, mon, flags, manager)
			}

			mu.Lock()
			count = 0
			mu.Unlock()
			log.Println("send metrics grpc end")

			select {
			case <-ctxGrSh.Done():
				tickerReport.Stop()
				tickerPool.Stop()
				shutdown = true
				log.Println("graceful shutdown signal received for grpc")
			default:
			}
		}
	}()

	//Условие завершения функции
	times := 0
	for !shutdown && (maxCount < 0 || times < maxCount) {
		times++
		time.Sleep(1 * time.Second)
	}
}

func retryToSendMetrics(c Clienter, ctx context.Context, metrics []*pb.Metric, encryptedMetrics []byte, tries []int) error {
	var err error
	for try := 0; try < len(tries); try++ {
		time.Sleep(time.Duration(tries[try]) * time.Second)
		err = c.SendMetrics(ctx, metrics, encryptedMetrics)
		if err == nil {
			break
		}
		logger.WriteErrorLog("Error in request by try:"+strconv.Itoa(try), err.Error())
	}
	return err
}

// SendMetricsByGRPC отправляет метрики(несколько раз в случае неудачной отправки)
func SendMetricsByGRPC(r request, c Clienter, monitor *storage.Monitor, flags *SystemConfigFlags, manager *crypto.Manager) {
	metrics := metricsHandler.GetMetricsCollection(monitor)
	pbMetrics := ConvertToProto(metrics)
	ctx, err := setHashByMetrics(r, pbMetrics, flags)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "SetHashByMetrics")
	}
	pbMetrics, encryptedMetrics, err := encryptMetrics(pbMetrics, manager)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "EncryptMetrics")
	}
	err = c.SendMetrics(ctx, pbMetrics, encryptedMetrics)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "Error in sending metrics")
		err = retryToSendMetrics(c, ctx, pbMetrics, encryptedMetrics, errors.TriesTimes)
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Error in sending metrics by retry")
		}
	}
}

// ConvertToProto преобразует ваши models.Metrics в protobuf Metric
func ConvertToProto(metrics []models.Metrics) []*pb.Metric {
	var pbMetrics []*pb.Metric

	for _, m := range metrics {
		pbMetric := &pb.Metric{
			Id:   m.ID,
			Type: pb.Metric_gauge,
		}

		switch m.MType {
		case "gauge":
			if m.Value != nil {
				pbMetric.Value = *m.Value
			}
		case "counter":
			if m.Delta != nil {
				pbMetric.Type = pb.Metric_counter
				pbMetric.Delta = *m.Delta
			}
		}
		pbMetrics = append(pbMetrics, pbMetric)
	}

	return pbMetrics
}
