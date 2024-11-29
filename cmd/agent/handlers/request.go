package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/ramil063/gometrics/cmd/agent/storage"
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
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
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
				if typeOfS.Field(i).Name == "PoolCount" {
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
	var m storage.Monitor
	url := "http://" + MainURL + "/updates"

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
			body := make([]byte, 100)
			v := reflect.ValueOf(m)
			typeOfS := v.Type()
			log.Println("send metrics json")
			allMetrics := make([]models.Metrics, 100)

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

			if err := c.SendPostRequestWithBody(url, body); err != nil {
				logger.WriteErrorLog("Error in request", err.Error())
			}
		}
	}
	return nil
}
