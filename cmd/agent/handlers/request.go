package handlers

import (
	"fmt"
	"github.com/ramil063/gometrics/cmd/agent/storage"
	"log"
	"net/http"
	"reflect"
	"time"
)

var MainUrl = "http://localhost:8080"

type RequestInterface interface {
	SendMetrics(c ClientInterface, maxCount int) error
}

type ClientInterface interface {
	SendPostRequest(url string) error
	NewRequest(method string, url string) (*http.Request, error)
}

type client struct {
	httpClient *http.Client
}

type request struct {
}

func NewClient() ClientInterface {
	return client{
		httpClient: &http.Client{},
	}
}

func NewRequest() RequestInterface {
	return request{}
}

func (c client) SendPostRequest(url string) error {
	req, _ := c.NewRequest("POST", url)
	req.Header.Set("Content-Type", "text/plain")
	_, err := c.httpClient.Do(req)
	return err
}

func (c client) NewRequest(method string, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

func (r request) SendMetrics(c ClientInterface, maxCount int) error {
	var pollInterval = time.Duration(storage.PollInterval) * time.Second
	var reportByPoll = int(storage.ReportInterval / storage.PollInterval)
	count := 0
	var m storage.Monitor

	for count < maxCount {
		<-time.After(pollInterval)
		m = storage.NewMonitor()
		m.PollCount = storage.Counter(count)
		log.Println("get metrics")
		count++

		if (count % reportByPoll) == 0 {
			v := reflect.ValueOf(m)
			typeOfS := v.Type()
			log.Println("send metrics")

			for i := 0; i < v.NumField(); i++ {
				metricType := "gauge"
				if typeOfS.Field(i).Name == "PoolCount" {
					metricType = "counter"
				}
				metricValue := fmt.Sprintf("%v", v.Field(i).Interface())
				url := MainUrl + "/update/" + metricType + "/" + typeOfS.Field(i).Name + "/" + metricValue

				err := c.SendPostRequest(url)
				if err != nil {
					log.Fatal("Error")
					return err
				}
			}
		}
	}
	return nil
}
