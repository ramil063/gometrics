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
	SendMultipleMetricsJSON(c JSONClienter, maxCount int) error
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
func (r request) SendMultipleMetricsJSON(c JSONClienter, maxCount int) error {
	var interval = 1 * time.Second
	count := 0
	seconds := 0
	url := "http://" + MainURL + "/updates"

	// создаем буферизованный канал для принятия метрик в сендер
	requestBodies := make(chan []byte)
	resultCh := make(chan Result)
	var wg sync.WaitGroup

	for maxCount < 0 || seconds < maxCount {
		<-time.After(interval)
		seconds++
		if (seconds % PollInterval) == 0 {
			select {
			case result := <-resultCh:
				if result.err != nil {
					logger.WriteErrorLog(result.err.Error(), "Collect or prepare metrics error")
					return result.err
				}
			default:
				wg.Add(1)
				monitor := CollectMonitorMetrics(&count)
				monitor = CollectGopsutilMetrics(monitor, resultCh)
				go CollectMetricsRequestBodies(requestBodies, monitor, &wg)
			}
		}

		if (seconds % ReportInterval) == 0 {
			select {
			case result := <-resultCh:
				if result.err != nil {
					logger.WriteErrorLog(result.err.Error(), "Send metrics error")
					return result.err
				}
			default:
				wg.Wait()
				for worker := 0; worker < RateLimit; worker++ {
					go SendMetrics(c, url, requestBodies, resultCh)
				}
				requestBodies = make(chan []byte)
			}
		}
	}
	return nil
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

func CollectMetricsRequestBodies(requestBodies chan []byte, monitor chan storage.Monitor, wg *sync.WaitGroup) {
	log.Println("get metrics json")
	wg.Done()

	m := <-monitor

	v := reflect.ValueOf(m)
	typeOfS := v.Type()
	allMetrics := make([]models.Metrics, 0, 100)

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
		allMetrics = append(allMetrics, metrics)
	}
	body, err := json.Marshal(allMetrics)
	if err != nil {
		logger.WriteErrorLog("Error marshal metrics", err.Error())
	}
	requestBodies <- body
}

func CollectMonitorMetrics(count *int) chan storage.Monitor {
	resultMonitor := make(chan storage.Monitor)

	go func() {
		defer close(resultMonitor)

		var m storage.Monitor
		m = storage.NewMonitor()
		m.PollCount = models.Counter(*count)
		*(count)++
		resultMonitor <- m
	}()

	return resultMonitor
}

func CollectGopsutilMetrics(monitors chan storage.Monitor, resultCh chan Result) chan storage.Monitor {
	resultMonitor := make(chan storage.Monitor)

	go func() {
		defer close(resultMonitor)

		gm, err := storage.NewGopsutilMonitor()
		if err != nil {
			result := Result{err: err}
			resultCh <- result
		}

		for m := range monitors {
			m.TotalMemory = gm.TotalMemory
			m.FreeMemory = gm.FreeMemory
			m.CPUutilization1 = gm.CPUutilization1
			resultMonitor <- m
		}
	}()

	return resultMonitor
}

func SendMetrics(c JSONClienter, url string, requestBodies chan []byte, resultCh chan Result) {
	log.Println("send metrics json")
	var err error

	for body := range requestBodies {
		if err = c.SendPostRequestWithBody(url, body); err != nil {
			logger.WriteErrorLog("Error in request", err.Error())
			var reqErr *internalErrors.RequestError
			if errors.Is(err, reqErr) || errors.Is(err, syscall.ECONNREFUSED) {
				err = retryToSendMetrics(c, url, body, internalErrors.TriesTimes)
				result := Result{err: err}
				resultCh <- result
				return
			}
		}
	}
}
