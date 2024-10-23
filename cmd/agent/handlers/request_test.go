package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RequestMock struct {
	Request request
}

type ClientMock struct {
	Client client
}
type JSONClientMock struct {
	ClientMock
}

func (c ClientMock) SendPostRequest(url string) error {
	_, err := c.NewRequest("POST", url)
	return err
}

func (c ClientMock) NewRequest(method string, url string) (*http.Request, error) {
	return httptest.NewRequest(method, url, nil), nil
}

func (c JSONClientMock) SendPostRequestWithBody(url string, body []byte) error {
	_ = httptest.NewRequest("POST", url, bytes.NewReader(body))
	return nil
}

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name string
		want Requester
	}{
		{"create new request", RequestMock{}.Request},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewRequest())
		})
	}
}

func Test_request_SendMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"send request"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{}
			assert.NoError(t, r.SendMetrics(ClientMock{}, 5))
		})
	}
}

func Test_request_SendMetricsJSON(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"send request"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{}
			assert.NoError(t, r.SendMetricsJSON(JSONClientMock{}, 5))
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"create new request", "handlers.client"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, reflect.ValueOf(NewClient()).Type().String(), "NewClient()")
		})
	}
}

func TestNewJSONClient(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"create new request", "handlers.client"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, reflect.ValueOf(NewJSONClient()).Type().String(), "NewJSONClient()")
		})
	}
}

func Test_client_SendPostRequest(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"send post", args{url: "http://" + MainURL + "/update/gauge/metric1/100"}, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ClientMock{}
			tt.wantErr(t, c.SendPostRequest(tt.args.url), fmt.Sprintf("SendPostRequest(%v)", tt.args.url))
		})
	}
}

func Test_client_SendPostRequestWithBody(t *testing.T) {
	type args struct {
		url  string
		body []byte
	}
	body, _ := json.Marshal(`{"id": "metric1", "type": "gauge", "value": 100.250}`)
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"send post", args{url: "http://" + MainURL + "/update", body: body}, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := JSONClientMock{}
			tt.wantErr(t, c.SendPostRequestWithBody(tt.args.url, tt.args.body), fmt.Sprintf("SendPostRequestWithBody(%v, %v)", tt.args.url, tt.args.body))
		})
	}
}
