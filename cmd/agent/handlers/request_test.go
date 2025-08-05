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
	var wg sync.WaitGroup
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
			wg.Add(1)
			r.SendMultipleMetricsJSON(JSONClientMock{}, 7, context.Background(), flags, &wg)
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
