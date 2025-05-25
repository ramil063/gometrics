package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/cmd/agent/storage"
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
		want Requester
		name string
	}{
		{
			want: RequestMock{}.Request,
			name: "create new request",
		},
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
			assert.NoError(t, r.SendMetricsJSON(JSONClientMock{}, 11))
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
		wantErr assert.ErrorAssertionFunc
		args    args
		name    string
	}{
		{
			name:    "send post",
			args:    args{url: "http://" + MainURL + "/update/gauge/metric1/100"},
			wantErr: assert.NoError,
		},
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
		wantErr assert.ErrorAssertionFunc
		name    string
		args    args
	}{
		{
			wantErr: assert.NoError,
			name:    "send post",
			args:    args{url: "http://" + MainURL + "/update", body: body},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := JSONClientMock{}
			tt.wantErr(t, c.SendPostRequestWithBody(tt.args.url, tt.args.body), fmt.Sprintf("SendPostRequestWithBody(%v, %v)", tt.args.url, tt.args.body))
		})
	}
}

func Test_request_SendMultipleMetricsJSON(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"send request"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{}
			r.SendMultipleMetricsJSON(JSONClientMock{}, 7)
		})
	}
}

func Test_compressData(t *testing.T) {
	type args struct {
		data []byte
	}
	want := []byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x4a, 0x4, 0x4, 0x0, 0x0, 0xff, 0xff, 0x43, 0xbe, 0xb7, 0xe8, 0x1, 0x0, 0x0, 0x0}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"test 1", args{data: []byte("a")}, want, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compressData(tt.args.data)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}

func BenchmarkCollectMetricsRequestBodies(b *testing.B) {
	var m storage.Monitor

	for i := 0; i < b.N; i++ {
		CollectMetricsRequestBodies(&m)
	}
}

func BenchmarkCollectMonitorMetrics(b *testing.B) {
	var m storage.Monitor
	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		CollectMonitorMetrics(1, &m, &wg)
	}
}

func BenchmarkCollectGopsutilMetrics(b *testing.B) {
	var m storage.Monitor
	var wg sync.WaitGroup
	storage.SetMetricsToMonitor(&m)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CollectGopsutilMetrics(&m, &wg)
	}
}
