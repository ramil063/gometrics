package handlers

import (
	"bytes"
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

type JSONRequester interface {
	SendMetricsJSON(c JSONClienter, maxCount int) error
	SendMultipleMetricsJSON(c JSONClienter, maxCount int)
}

type Requester interface {
	JSONRequester
	SendMetrics(c Clienter, maxCount int) error
}

type Clienter interface {
	SendPostRequest(url string) error
	NewRequest(method string, url string) (*http.Request, error)
}

type JSONClienter interface {
	Clienter
	SendPostRequestWithBody(url string, body []byte) error
}

type client struct {
	httpClient *http.Client
}

type request struct {
}

// структура, в которую добавили ошибку
type Result struct {
	err error
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

// SendPostRequestWithBody отправка пост запроса с телом
func (c client) SendPostRequestWithBody(url string, body []byte) error {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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
	var m storage.Monitor

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
			v := reflect.ValueOf(m)
			typeOfS := v.Type()
			log.Println("send metrics")

			for i := 0; i < v.NumField(); i++ {
				metricType := "gauge"
				if typeOfS.Field(i).Name == "PollCount" {
					metricType = "counter"
				}
				metricValue := fmt.Sprintf("%v", v.Field(i).Interface())
				url := "http://" + MainURL + "/update/" + metricType + "/" + typeOfS.Field(i).Name + "/" + metricValue

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
	var m storage.Monitor

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
			v := reflect.ValueOf(m)
			typeOfS := v.Type()
			log.Println("send metrics json")

			for i := 0; i < v.NumField(); i++ {
				metricID := typeOfS.Field(i).Name
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
	times := 0
	url := "http://" + MainURL + "/updates"

	// создаем канал для принятия метрик в сендер
	requestBodies := make(chan []byte)
	monitors := make(chan storage.Monitor)

	tickerPool := time.NewTicker(pollInterval)
	defer tickerPool.Stop()
	tickerReport := time.NewTicker(reportInterval)
	defer tickerReport.Stop()

	var sendWg sync.WaitGroup
	for maxCount < 0 || times < maxCount {
		times++
		select {
		case <-tickerPool.C:
			var wg sync.WaitGroup
			wg.Add(2)
			monitor := CollectMonitorMetrics(&count, &wg)
			CollectGopsutilMetrics(monitors, monitor, &wg)
			log.Println("get metrics json")
			wg.Wait()
			sendWg.Add(1)
		case <-tickerReport.C:
			go CollectMetricsRequestBodies(requestBodies, monitors, &sendWg)

			for worker := 0; worker < RateLimit; worker++ {
				go SendMetrics(c, url, requestBodies, &sendWg)
			}
			requestBodies = make(chan []byte)
			monitors = make(chan storage.Monitor)
		}
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

func CollectMetricsRequestBodies(requestBodies chan []byte, monitors chan storage.Monitor, wg *sync.WaitGroup) {
	defer close(monitors)

	for m := range monitors {
		v := reflect.ValueOf(m)
		typeOfS := v.Type()
		allMetrics := make([]models.Metrics, 0, 100)

		for i := 0; i < v.NumField(); i++ {
			metricID := typeOfS.Field(i).Name
			// Пропускаем поле "mx"
			if metricID == "mx" {
				continue
			}

			metricValue, _ := strconv.ParseFloat(fmt.Sprintf("%v", v.Field(i).Interface()), 64)
			delta := int64(m.PollCount)

			if metricID == "CPUutilization" {
				CPUutilization := m.GetAllCPUutilization()
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
		requestBodies <- body
		wg.Done()
	}
}

func CollectMonitorMetrics(count *int, wg *sync.WaitGroup) chan storage.Monitor {
	resultMonitor := make(chan storage.Monitor)
	defer wg.Done()

	go func() {
		defer close(resultMonitor)

		m := storage.NewMonitor()
		m.PollCount = models.Counter(*count)
		*(count)++
		resultMonitor <- m
	}()

	return resultMonitor
}

func CollectGopsutilMetrics(monitors chan storage.Monitor, monitor chan storage.Monitor, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		gm, err := storage.NewGopsutilMonitor()
		if err != nil {
			logger.WriteErrorLog(err.Error(), "NewGopsutilMonitor")
		}
		m := <-monitor
		m.TotalMemory = gm.TotalMemory
		m.FreeMemory = gm.FreeMemory
		CPUutilization := gm.GetAllCPUutilization()
		for key, value := range CPUutilization {
			m.StoreCPUutilizationValue(key, value)
		}
		monitors <- m
	}()
}

func SendMetrics(c JSONClienter, url string, requestBodies chan []byte, wg *sync.WaitGroup) {
	defer wg.Wait()
	log.Println("send metrics json")
	var err error

	for body := range requestBodies {
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
}
