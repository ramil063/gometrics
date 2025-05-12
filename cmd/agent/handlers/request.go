package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ramil063/gometrics/cmd/agent/storage"
	internalErrors "github.com/ramil063/gometrics/internal/errors"
	"github.com/ramil063/gometrics/internal/hash"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

// JSONRequester отправляет данные в формате json
type JSONRequester interface {
	SendMetricsJSON(c JSONClienter, maxCount int) error
	SendMultipleMetricsJSON(c JSONClienter, maxCount int)
}

// Requester отправляет данные
type Requester interface {
	JSONRequester
	SendMetrics(c Clienter, maxCount int) error
}

// Clienter работа с клиентов
type Clienter interface {
	SendPostRequest(url string) error
	NewRequest(method string, url string) (*http.Request, error)
}

// JSONClienter работа с клиентом в формате json
type JSONClienter interface {
	Clienter
	SendPostRequestWithBody(url string, body []byte) error
}

type client struct {
	httpClient *http.Client
}

type request struct {
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
	return request{}
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
func (c client) SendPostRequestWithBody(url string, body []byte) error {
	data, err := compressData(body)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	if HashKey != "" {
		hashSha256 := hash.CreateSha256(body, HashKey)
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
func (r request) SendMetrics(c Clienter, maxCount int) error {
	var interval = 1 * time.Second
	count := 0
	seconds := 0
	var m *storage.Monitor

	for count < maxCount {
		<-time.After(interval)
		seconds++
		if (seconds % PollInterval) == 0 {
			log.Println("get metrics")
			m = storage.NewMonitor()
			m.PollCount = models.Counter(count)
			count++
		}

		if (seconds % ReportInterval) == 0 {
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
				url := "http://" + MainURL + "/update/" + metricType + "/" + name + "/" + metricValue

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
func (r request) SendMetricsJSON(c JSONClienter, maxCount int) error {
	var interval = 1 * time.Second
	count := 0
	seconds := 0
	var m *storage.Monitor

	for seconds < maxCount {
		<-time.After(interval)
		seconds++
		if (seconds % PollInterval) == 0 {
			log.Println("get metrics json")
			m = storage.NewMonitor()
			m.PollCount = models.Counter(count)
			count++
		}

		if (seconds % ReportInterval) == 0 {
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

				url := "http://" + MainURL + "/update"
				body, err := json.Marshal(metrics)
				if err != nil {
					logger.WriteErrorLog("Error marshal metrics", err.Error())
				}

				err = c.SendPostRequestWithBody(url, body)
				if err != nil {
					logger.WriteErrorLog("Error in request", err.Error())
				}
			}
		}
	}
	return nil
}

// SendMultipleMetricsJSON отправка нескольких метрик
func (r request) SendMultipleMetricsJSON(c JSONClienter, maxCount int) {
	var pollInterval = time.Duration(PollInterval) * time.Second
	var reportInterval = time.Duration(ReportInterval) * time.Second
	count := 0
	url := "http://" + MainURL + "/updates"

	var monitor storage.Monitor
	tickerPool := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)

	var sendMonitor = make(chan *storage.Monitor, 1)
	defer close(sendMonitor)

	var mu sync.Mutex

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
			CollectMonitorMetrics(count, &monitor, &collectWg)
			CollectGopsutilMetrics(&monitor, &collectWg)
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

			for worker := 0; worker < RateLimit; worker++ {
				go SendMetrics(c, url, mon, worker)
			}

			mu.Lock()
			count = 0
			mu.Unlock()
			log.Println("send metrics json end")
		}
	}()

	//Условие завершения функции(для тестирования)
	times := 0
	for maxCount < 0 || times < maxCount {
		times++
		time.Sleep(1 * time.Second)
	}
}

func retryToSendMetrics(c JSONClienter, url string, body []byte, tries []int) error {
	var err error
	for try := 0; try < len(tries); try++ {
		time.Sleep(time.Duration(tries[try]) * time.Second)
		err = c.SendPostRequestWithBody(url, body)
		if err == nil {
			break
		}
		logger.WriteErrorLog("Error in request by try:"+strconv.Itoa(try), err.Error())
	}
	return err
}

// CollectMetricsRequestBodies сбор метрик в тела для отправки на сторонний сервис
func CollectMetricsRequestBodies(monitor *storage.Monitor) []byte {
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
	body, err := json.Marshal(allMetrics)
	if err != nil {
		logger.WriteErrorLog("Error marshal metrics", err.Error())
	}
	return body
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

// SendMetrics отправляет метрики(несколько раз в случае неудачной отправки)
func SendMetrics(c JSONClienter, url string, monitor *storage.Monitor, worker int) {
	var err error
	body := CollectMetricsRequestBodies(monitor)

	if err = c.SendPostRequestWithBody(url, body); err != nil {
		logger.WriteErrorLog(err.Error(), "Error in request")
		var reqErr *internalErrors.RequestError
		if errors.Is(err, reqErr) || errors.Is(err, syscall.ECONNREFUSED) {
			err = retryToSendMetrics(c, url, body, internalErrors.TriesTimes)
			if err != nil {
				logger.WriteErrorLog(err.Error(), "Error in retry request")
			}
		}
	}
}

// compressData Функция для сжатия данных
func compressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	if err = gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
