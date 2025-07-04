package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/models"
)

func Test_update(t *testing.T) {
	handlers.Restore = false
	tests := []struct {
		name       string
		statusCode int
	}{
		{"test 1", 200},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/gauge/a/1", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			updateHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
				ms := NewMemStorage()
				Update(rw, req, ms)
			}
			handlerToTest := http.HandlerFunc(updateHandlerFunction)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.statusCode, res.StatusCode)
			defer res.Body.Close()
			// получаем и проверяем тело запроса
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	req.Header.Set("Content-Type", "text/plain")
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	handlers.Restore = false
	ms := NewMemStorage()
	ts := httptest.NewServer(Router(ms))
	defer ts.Close()

	var testTable = []struct {
		name   string
		url    string
		status int
	}{
		{"test 1", "/update/gauge/a/1", http.StatusOK},
		{"test 2", "/update/gauge/a", http.StatusBadRequest},
		{"test 3", "/update/counter/", http.StatusNotFound},
		{"test 4", "/update/counter/testSetGet32/417", http.StatusOK},
	}
	for _, v := range testTable {
		resp, _ := testRequest(t, ts, "POST", v.url)
		defer resp.Body.Close()
		assert.Equal(t, v.status, resp.StatusCode)
	}
}

func Test_getValue(t *testing.T) {
	handlers.Restore = false
	ms := NewMemStorage()
	ts := httptest.NewServer(Router(ms))
	defer ts.Close()

	type want struct {
		response string
		code     int
	}
	testsG := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "gauge test 1",
			url:  "/value/gauge/a",
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name: "gauge test 2",
			url:  "/value/gauge/a1",
			want: want{
				code:     404,
				response: "",
			},
		},
		{
			name: "gauge test 3",
			url:  "/value/gauge/",
			want: want{
				code:     404,
				response: "",
			},
		},
		{
			name: "gauge test 4",
			url:  "/value/gauge/testUnknown80",
			want: want{
				code:     404,
				response: "",
			},
		},
	}
	for _, test := range testsG {
		t.Run(test.name, func(t *testing.T) {
			err := ms.SetGauge("a", 1.1)
			assert.NoError(t, err)
			resp, _ := testRequest(t, ts, "GET", test.url)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
		})
	}
	testsC := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "counter test 1",
			url:  "/value/counter/a",
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name: "counter test 2",
			url:  "/value/counter/a1",
			want: want{
				code:     404,
				response: "",
			},
		},
		{
			name: "counter test 3",
			url:  "/value/counter/",
			want: want{
				code:     404,
				response: "",
			},
		},
	}
	for _, test := range testsC {
		t.Run(test.name, func(t *testing.T) {
			err := ms.AddCounter("a", 1)
			assert.NoError(t, err)
			resp, _ := testRequest(t, ts, "GET", test.url)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
		})
	}
}

func Test_home(t *testing.T) {
	handlers.Restore = false
	ms := NewMemStorage()
	ts := httptest.NewServer(Router(ms))
	defer ts.Close()

	type want struct {
		response    string
		contentType string
		code        int
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "test 1",
			url:  "/",
			want: want{
				code:        200,
				response:    "",
				contentType: "text/html",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ms.SetGauge("a", 1.1)
			assert.NoError(t, err)
			resp, _ := testRequest(t, ts, "GET", test.url)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test 1", "*memory.MemStorage"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStorage()
			assert.Equalf(t, tt.want, reflect.ValueOf(ms).Type().String(), "NewMemStorage()")
		})
	}
}

func Test_updateMetricsJSON(t *testing.T) {
	handlers.Restore = false
	updateMetricsJSONHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
		UpdateMetricsJSON(rw, req, NewMemStorage())
	}
	handler := http.HandlerFunc(updateMetricsJSONHandlerFunction)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string // добавляем название тестов
		method       string
		body         string // добавляем тело запроса в табличные тесты
		expectedBody string
		expectedCode int
	}{
		{
			name:         "test 1",
			method:       http.MethodPost,
			expectedBody: "",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "test 2",
			method:       http.MethodPost,
			body:         `{"id": "met1", "type": "gauge", "value": 1.1}`,
			expectedBody: `{"id": "met1", "type": "gauge", "value": 1.1}`,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL

			if len(tc.body) > 0 {
				req.SetHeader("Content-Type", "application/json")
				req.SetBody(tc.body)
			}

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
			}
		})
	}
}

func Test_getValueMetricsJSON(t *testing.T) {
	filePath := "../../../../internal/storage/files/test.json"

	handlers.FileStoragePath = filePath
	getValueMetricsJSONHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
		s := GetStorage("", "")
		_ = s.SetGauge("met1", 1.1)
		GetValueMetricsJSON(rw, req, s)
	}
	handler := http.HandlerFunc(getValueMetricsJSONHandlerFunction)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string // добавляем название тестов
		method       string
		body         string // добавляем тело запроса в табличные тесты
		expectedBody string
		expectedCode int
	}{
		{
			name:         "test 1",
			method:       http.MethodPost,
			expectedBody: "",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "test 2",
			method:       http.MethodPost,
			body:         `{"id": "met1", "type": "gauge", "value":1.1}`,
			expectedBody: `{"id": "met1", "type": "gauge", "value":1.1}`,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL

			if len(tc.body) > 0 {
				req.SetHeader("Content-Type", "application/json")
				req.SetBody(tc.body)
			}

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
			}
		})
	}
}

func Test_updates(t *testing.T) {
	var mock sqlmock.Sqlmock
	dml.DBRepository.Database, mock, _ = sqlmock.New()
	defer dml.DBRepository.Database.Close()

	mock.ExpectExec("^INSERT INTO gauge *").
		WithArgs("met1", float64(1.1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^INSERT INTO gauge *").
		WithArgs("met2", float64(2.2)).
		WillReturnResult(sqlmock.NewResult(2, 1))

	updatesHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
		s := &db.Storage{}
		Updates(rw, req, s)
	}
	handler := http.HandlerFunc(updatesHandlerFunction)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string // добавляем название тестов
		method       string
		body         string // добавляем тело запроса в табличные тесты
		expectedBody string
		expectedCode int
	}{
		{
			name:         "test 1",
			method:       http.MethodPost,
			expectedBody: "",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "test 2",
			method:       http.MethodPost,
			body:         `[{"id": "met1", "type": "gauge", "value":1.1},{"id": "met2", "type": "gauge", "value":2.2}]`,
			expectedCode: http.StatusOK,
			expectedBody: `[{"id": "met1", "type": "gauge", "value":1.1},{"id": "met2", "type": "gauge", "value":2.2}]`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL

			if len(tc.body) > 0 {
				req.SetHeader("Content-Type", "application/json")
				req.SetBody(tc.body)
			}

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
			}
		})
	}
}

func Test_updateMetrics(t *testing.T) {
	dbs := NewMemStorage()

	tests := []struct {
		name    string
		metrics []models.Metrics
		want    []models.Metrics
		delta   int64
	}{
		{
			name: "test 1",
			metrics: []models.Metrics{
				{
					ID:    "met1",
					MType: "counter",
				},
			},
			want: []models.Metrics{
				{
					ID:    "met1",
					MType: "counter",
				},
			},
			delta: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metrics[0].Delta = &tt.delta

			got, err := UpdateMetrics(dbs, tt.metrics)
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.want[0].ID, got[0].ID)
			assert.Equal(t, tt.want[0].MType, got[0].MType)
			assert.Equal(t, tt.delta, *(got[0].Delta))
		})
	}
}

func BenchmarkUpdateMetrics(b *testing.B) {
	dbs := NewMemStorage()
	metrics := []models.Metrics{
		{
			ID:    "met1",
			MType: "counter",
		},
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = UpdateMetrics(dbs, metrics)
	}
}
