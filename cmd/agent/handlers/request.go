package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ramil063/gometrics/cmd/agent/handlers/gzip"
	metricsHandler "github.com/ramil063/gometrics/cmd/agent/handlers/metrics"
	"github.com/ramil063/gometrics/cmd/agent/storage"
	internalErrors "github.com/ramil063/gometrics/internal/errors"
	"github.com/ramil063/gometrics/internal/hash"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

// JSONRequester отправляет данные в формате json
type JSONRequester interface {
	SendMetricsJSON(c JSONClienter, maxCount int, flags *SystemConfigFlags, manager *crypto.Manager) error
	SendMultipleMetricsJSON(c JSONClienter, maxCount int, ctxGrSh context.Context, flags *SystemConfigFlags, manager *crypto.Manager, serversWg *sync.WaitGroup)
}

// Requester отправляет данные
type Requester interface {
	JSONRequester
	SendMetrics(c Clienter, maxCount int, flags *SystemConfigFlags) error
}

// Clienter работа с клиентов
type Clienter interface {
	SendPostRequest(url string) error
	NewRequest(method string, url string) (*http.Request, error)
}

// JSONClienter работа с клиентом в формате json
type JSONClienter interface {
	Clienter
	SendPostRequestWithBody(r request, url string, body []byte, flags *SystemConfigFlags, manager *crypto.Manager) error
}

type client struct {
	httpClient *http.Client
}

type request struct {
	IP string
}

func NewClient() Clienter {
	return client{
		httpClient: &http.Client{},
	}
}

func NewJSONClient() JSONClienter {
	return client{
		httpClient: &http.Client{},
	}
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

// SendPostRequest отправка пост запроса
func (c client) SendPostRequest(url string) error {
	req, _ := c.NewRequest("POST", url)
	req.Header.Set("Content-Type", "text/plain")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// SendPostRequestWithBody отправляет пост запроса с телом
func (c client) SendPostRequestWithBody(r request, url string, body []byte, flags *SystemConfigFlags, manager *crypto.Manager) error {
	var err error
	data := body

	encryptor := manager.GetDefaultEncryptor()
	if encryptor != nil {
		data, err = encryptor.Encrypt(data)
		if err != nil {
			return err
		}
	}

	data, err = gzip.CompressData(data)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("X-Real-IP", r.IP)

	if flags.HashKey != "" {
		hashSha256 := hash.CreateSha256(body, flags.HashKey)
		req.Header.Set("HashSHA256", hashSha256)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var b []byte
	res.Body.Read(b)
	if res.StatusCode != http.StatusOK {
		return internalErrors.NewRequestError(res.Status, res.StatusCode)
	}

	return nil
}

// NewRequest создание нового запроса
func (c client) NewRequest(method string, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

// SendMetrics отправка метрик
func (r request) SendMetrics(c Clienter, maxCount int, flags *SystemConfigFlags) error {
	var interval = 1 * time.Second
	count := 0
	seconds := 0
	var m *storage.Monitor

	for count < maxCount {
		<-time.After(interval)
		seconds++
		if (seconds % flags.PollInterval) == 0 {
			log.Println("get metrics")
			m = storage.NewMonitor()
			m.PollCount = models.Counter(count)
			count++
		}

		if (seconds % flags.ReportInterval) == 0 {
			v := reflect.ValueOf(m).Elem()
			typeOfS := v.Type()
			log.Println("send metrics")

			for i := 0; i < v.NumField(); i++ {
				metricType := "gauge"
				name := typeOfS.Field(i).Name
				if name == "mx" {
					continue
				}
				if name == "PollCount" {
					metricType = "counter"
				}
				metricValue := fmt.Sprintf("%v", v.Field(i).Interface())
				url := "http://" + flags.Address + "/update/" + metricType + "/" + name + "/" + metricValue

				err := c.SendPostRequest(url)
				if err != nil {
					logger.WriteErrorLog("Send request error", err.Error())
					log.Fatal("Error", err)
					return err
				}
			}
		}
	}
	return nil
}

// SendMetricsJSON отправка метрик
func (r request) SendMetricsJSON(c JSONClienter, maxCount int, flags *SystemConfigFlags, manager *crypto.Manager) error {
	var interval = 1 * time.Second
	count := 0
	seconds := 0
	var m *storage.Monitor

	for seconds < maxCount {
		<-time.After(interval)
		seconds++
		if (seconds % flags.PollInterval) == 0 {
			log.Println("get metrics json")
			m = storage.NewMonitor()
			m.PollCount = models.Counter(count)
			count++
		}

		if (seconds % flags.ReportInterval) == 0 {
			v := reflect.ValueOf(m).Elem()
			typeOfS := v.Type()
			log.Println("send metrics json")

			for i := 0; i < v.NumField(); i++ {
				metricID := typeOfS.Field(i).Name
				if metricID == "mx" {
					continue
				}

				metricValue, _ := strconv.ParseFloat(fmt.Sprintf("%v", v.Field(i).Interface()), 64)
				delta := int64(m.PollCount)

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

				url := "http://" + flags.Address + "/update"
				body, err := json.Marshal(metrics)
				if err != nil {
					logger.WriteErrorLog("Error marshal metrics", err.Error())
				}

				err = c.SendPostRequestWithBody(r, url, body, flags, manager)
				if err != nil {
					logger.WriteErrorLog("Error in request", err.Error())
				}
			}
		}
	}
	return nil
}

// SendMultipleMetricsJSON отправка нескольких метрик
func (r request) SendMultipleMetricsJSON(
	c JSONClienter,
	maxCount int,
	ctxGrSh context.Context,
	flags *SystemConfigFlags,
	manager *crypto.Manager,
	serversWg *sync.WaitGroup,
) {
	defer serversWg.Done()
	var pollInterval = time.Duration(flags.PollInterval) * time.Second
	var reportInterval = time.Duration(flags.ReportInterval) * time.Second
	count := 0
	url := "http://" + flags.Address + "/updates"

	var monitor storage.Monitor
	tickerPool := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)

	var sendMonitor = make(chan *storage.Monitor, 1)
	defer close(sendMonitor)

	var mu sync.Mutex
	var shutdown = false

	log.Println("agent start")

	go func() {
		defer tickerPool.Stop()
		for maxCount < 0 {
			<-tickerPool.C

			mu.Lock()
			count++
			mu.Unlock()

			log.Println("get metrics json start")
			var collectWg sync.WaitGroup
			metricsHandler.CollectMonitorMetrics(count, &monitor, &collectWg)
			metricsHandler.CollectGopsutilMetrics(&monitor, &collectWg)
			collectWg.Wait()

			sendMonitor <- &monitor
			log.Println("get metrics json end")
			if len(sendMonitor) == 1 {
				<-sendMonitor
			}
		}
	}()

	go func() {
		defer tickerReport.Stop()
		for maxCount < 0 {
			<-tickerReport.C

			log.Println("send metrics json start")
			mon := <-sendMonitor
			log.Println("send metrics json count value=", mon.GetCountValue())

			for worker := 0; worker < flags.RateLimit; worker++ {
				go SendMetrics(r, c, url, mon, flags, manager)
			}

			mu.Lock()
			count = 0
			mu.Unlock()
			log.Println("send metrics json end")

			select {
			case <-ctxGrSh.Done():
				tickerReport.Stop()
				tickerPool.Stop()
				shutdown = true
				log.Println("graceful shutdown signal received")
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

func retryToSendMetrics(r request, c JSONClienter, url string, body []byte, tries []int, flags *SystemConfigFlags, manager *crypto.Manager) error {
	var err error
	for try := 0; try < len(tries); try++ {
		time.Sleep(time.Duration(tries[try]) * time.Second)
		err = c.SendPostRequestWithBody(r, url, body, flags, manager)
		if err == nil {
			break
		}
		logger.WriteErrorLog("Error in request by try:"+strconv.Itoa(try), err.Error())
	}
	return err
}

// SendMetrics отправляет метрики(несколько раз в случае неудачной отправки)
func SendMetrics(r request, c JSONClienter, url string, monitor *storage.Monitor, flags *SystemConfigFlags, manager *crypto.Manager) {
	var err error
	body := metricsHandler.CollectMetricsRequestBodies(monitor)

	if err = c.SendPostRequestWithBody(r, url, body, flags, manager); err != nil {
		logger.WriteErrorLog(err.Error(), "Error in request")
		var reqErr *internalErrors.RequestError
		if errors.Is(err, reqErr) || errors.Is(err, syscall.ECONNREFUSED) {
			err = retryToSendMetrics(r, c, url, body, internalErrors.TriesTimes, flags, manager)
			if err != nil {
				logger.WriteErrorLog(err.Error(), "Error in retry request")
			}
		}
	}
}
