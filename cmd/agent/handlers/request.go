package handlers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/ramil063/gometrics/cmd/agent/storage"
)

type Requester interface {
	SendMetrics(c Clienter, maxCount int) error
}

type Clienter interface {
	SendPostRequest(url string) error
	NewRequest(method string, url string) (*http.Request, error)
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

func NewRequest() Requester {
	return request{}
}

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

func (c client) NewRequest(method string, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

func (r request) SendMetrics(c Clienter, maxCount int) error {
	var interval = 1 * time.Second
	count := 0
	seconds := 0
	var m storage.Monitor

	for count < maxCount {
		<-time.After(interval)
		seconds++
		if (seconds % int(PollInterval)) == 0 {
			log.Println("get metrics")
			m = storage.NewMonitor()
			m.PollCount = storage.Counter(count)
			count++
		}

		if (seconds % int(ReportInterval)) == 0 {
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
					log.Fatal("Error", err)
					return err
				}
			}
		}
	}
	return nil
}
