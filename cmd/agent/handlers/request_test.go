package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/cmd/agent/storage"
	internalErrors "github.com/ramil063/gometrics/internal/errors"
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

func (c JSONClientMock) SendPostRequestWithBody(r request, url string, body []byte, flags *SystemConfigFlags) error {
	_ = httptest.NewRequest("POST", url, bytes.NewReader(body))
	return nil
}

func TestNewRequest(t *testing.T) {
	req := RequestMock{}
	req.Request = request{
		IP: "127.0.1.1",
	}
	tests := []struct {
		want Requester
		name string
	}{
		{
			want: req.Request,
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

	flags := &SystemConfigFlags{
		Address:        ":8080",
		PollInterval:   2,
		ReportInterval: 10,
		RateLimit:      1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{}
			assert.NoError(t, r.SendMetrics(ClientMock{}, 5, flags))
		})
	}
}

func Test_request_SendMetricsJSON(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"send request"},
	}
	flags := &SystemConfigFlags{
		Address:        ":8080",
		PollInterval:   2,
		ReportInterval: 10,
		RateLimit:      1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{}
			assert.NoError(t, r.SendMetricsJSON(JSONClientMock{}, 11, flags))
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
			args:    args{url: "http://localhost:8080/update/gauge/metric1/100"},
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
			args:    args{url: "http://localhost:8080/update", body: body},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := JSONClientMock{}
			flags := SystemConfigFlags{}
			tt.wantErr(
				t,
				c.SendPostRequestWithBody(request{}, tt.args.url, tt.args.body, &flags),
				fmt.Sprintf("SendPostRequestWithBody(%v, %v)", tt.args.url, tt.args.body),
			)
		})
	}
}

func Test_request_SendMultipleMetricsJSON(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"send request"},
	}
	flags := &SystemConfigFlags{
		Address:        ":8080",
		PollInterval:   2,
		ReportInterval: 10,
		RateLimit:      1,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{}
			r.SendMultipleMetricsJSON(JSONClientMock{}, 7, context.Background(), flags)
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

func TestCollectMetricsRequestBodies(t *testing.T) {
	tests := []struct {
		name    string
		monitor *storage.Monitor
		want    []byte
	}{
		{
			name:    "test 1",
			monitor: &storage.Monitor{},
			want:    []byte("[{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":0},{\"id\":\"BuckHashSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"Frees\",\"type\":\"gauge\",\"value\":0},{\"id\":\"GCSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapAlloc\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapIdle\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapInuse\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapObjects\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapReleased\",\"type\":\"gauge\",\"value\":0},{\"id\":\"HeapSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"LastGC\",\"type\":\"gauge\",\"value\":0},{\"id\":\"Lookups\",\"type\":\"gauge\",\"value\":0},{\"id\":\"MCacheInuse\",\"type\":\"gauge\",\"value\":0},{\"id\":\"MCacheSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"Mallocs\",\"type\":\"gauge\",\"value\":0},{\"id\":\"NextGC\",\"type\":\"gauge\",\"value\":0},{\"id\":\"OtherSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"PauseTotalNs\",\"type\":\"gauge\",\"value\":0},{\"id\":\"StackInuse\",\"type\":\"gauge\",\"value\":0},{\"id\":\"Sys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"StackSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"MSpanInuse\",\"type\":\"gauge\",\"value\":0},{\"id\":\"MSpanSys\",\"type\":\"gauge\",\"value\":0},{\"id\":\"TotalMemory\",\"type\":\"gauge\",\"value\":0},{\"id\":\"FreeMemory\",\"type\":\"gauge\",\"value\":0},{\"id\":\"TotalAlloc\",\"type\":\"gauge\",\"value\":0},{\"id\":\"GCCPUFraction\",\"type\":\"gauge\",\"value\":0},{\"id\":\"NumForcedGC\",\"type\":\"gauge\",\"value\":0},{\"id\":\"NumGC\",\"type\":\"gauge\",\"value\":0},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":0},{\"id\":\"RandomValue\",\"type\":\"gauge\",\"value\":0}]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodies := CollectMetricsRequestBodies(tt.monitor)
			assert.Equal(t, tt.want, bodies)
		})
	}
}

func TestSendPostRequest_Success(t *testing.T) {
	// Создаем мок-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод и заголовок
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			t.Error("Expected Content-Type: text/plain")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := client{httpClient: &http.Client{}}
	err := c.SendPostRequest(ts.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func Test_client_SendPostRequestWithBody1(t *testing.T) {
	// Создаем мок-сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод и заголовок
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Expected Content-Type: application/json")
		}
		if r.Header.Get("Content-Encoding") != "gzip" {
			t.Error("Expected Content-Encoding: gzip")
		}
		if r.Header.Get("Accept-Encoding") != "gzip" {
			t.Error("Expected Accept-Encoding: gzip")
		}
		if r.Header.Get("X-Real-IP") != "127.0.0.1" {
			t.Error(r.Header.Get("X-Real-IP"))
			t.Error("Expected X-Real-IP: gzip")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := client{httpClient: &http.Client{}}
	r := request{
		IP: "127.0.0.1",
	}
	flags := SystemConfigFlags{}
	err := c.SendPostRequestWithBody(r, ts.URL, []byte("a"), &flags)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

type MockClient struct {
	ExpectedURL  string
	ExpectedBody []byte
	Attempts     int
	ShouldFail   bool
}

func (m *MockClient) SendPostRequestWithBody(r request, url string, body []byte, flags *SystemConfigFlags) error {
	m.Attempts++

	if url != m.ExpectedURL {
		return fmt.Errorf("unexpected URL: %s", url)
	}

	if !bytes.Equal(body, m.ExpectedBody) {
		return fmt.Errorf("unexpected body")
	}

	if m.ShouldFail {
		return &internalErrors.RequestError{Err: errors.New("mock error")}
	}
	return nil
}

func (m *MockClient) SendPostRequest(url string) error {
	return nil
}

func (m *MockClient) NewRequest(method string, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

func TestSendMetrics_SuccessWithMock(t *testing.T) {
	// 1. Подготавливаем мок
	mockClient := &MockClient{
		ExpectedURL:  "http://test",
		ExpectedBody: []byte(`{"metric":"value"}`), // Зависит от вашего CollectMetricsRequestBodies
	}

	// 2. Тестовые данные
	r := request{}
	monitor := &storage.Monitor{} // Заполните данными, которые вернут ожидаемый body
	flags := &SystemConfigFlags{}

	// 3. Запуск
	SendMetrics(r, mockClient, mockClient.ExpectedURL, monitor, flags)

	// 4. Проверки
	if mockClient.Attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", mockClient.Attempts)
	}
}
