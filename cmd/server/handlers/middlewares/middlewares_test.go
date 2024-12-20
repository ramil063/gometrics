package middlewares

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ramil063/gometrics/cmd/server/storage/memory"
	"github.com/ramil063/gometrics/internal/models"
)

func TestCheckMethodMw(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name        string
		method      string
		contentType string
		want        want
	}{
		{"test 1", http.MethodGet, "text-plain", want{200, "text/plain; charset=utf-8"}},
		{"test 2", http.MethodPost, "application/json", want{200, "text/plain; charset=utf-8"}},
		{"test 3", http.MethodPost, "text/plain", want{200, "text/plain; charset=utf-8"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/update/gauge/a/1", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", test.contentType)
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			})

			handlerToTest := CheckMethodMw(nextHandler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestCheckUpdateMetricsValue(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name       string
		pathValues map[string]string
		want       want
	}{
		{"test 1", map[string]string{"type": "gauge", "metric": "a", "value": "1"}, want{http.StatusOK, "", "text/plain; charset=utf-8"}},
		{"test 2", map[string]string{"type": "gauge", "metric": "a"}, want{http.StatusBadRequest, "", ""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			request.SetPathValue("type", test.pathValues["type"])
			if metric, ok := test.pathValues["metric"]; ok {
				request.SetPathValue("metric", metric)
			}
			if value, ok := test.pathValues["value"]; ok {
				request.SetPathValue("value", value)
			}
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			})

			handlerToTest := CheckUpdateMetricsValueMw(nextHandler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestCheckValueMetricsMw(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name       string
		pathValues map[string]string
		want       want
	}{
		{"test 1", map[string]string{"type": "gauge", "metric": "a"}, want{http.StatusOK, "", "text/plain; charset=utf-8"}},
		{"test 2", map[string]string{"type": "gauge1", "metric": "a"}, want{http.StatusBadRequest, "", ""}},
		{"test 4", map[string]string{"type": "counter"}, want{http.StatusNotFound, "", ""}},
		{"test 4", map[string]string{"type": "counter", "metric": "a"}, want{http.StatusOK, "", "text/plain; charset=utf-8"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/value", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			request.SetPathValue("type", test.pathValues["type"])
			if metric, ok := test.pathValues["metric"]; ok {
				request.SetPathValue("metric", metric)
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			})

			handlerToTest := CheckValueMetricsMw(nextHandler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestCheckUpdateMetricsNameMw(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name       string
		pathValues map[string]string
		want       want
	}{
		{"test 1", map[string]string{"type": "gauge", "metric": "a"}, want{http.StatusOK, "", "text/plain; charset=utf-8"}},
		{"test 4", map[string]string{"type": "counter"}, want{http.StatusBadRequest, "", ""}},
		{"test 4", map[string]string{"type": "counter", "metric": "a"}, want{http.StatusOK, "", "text/plain; charset=utf-8"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/value", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			request.SetPathValue("type", test.pathValues["type"])
			if metric, ok := test.pathValues["metric"]; ok {
				request.SetPathValue("metric", metric)
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			})

			handlerToTest := CheckUpdateMetricsNameMw(nextHandler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestCheckMetricsTypeMw(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name       string
		pathValues map[string]string
		want       want
	}{
		{"test 1", map[string]string{"type": "gauge", "metric": "a"}, want{http.StatusOK, "", "text/plain; charset=utf-8"}},
		{"test 4", map[string]string{"type": "counter1", "metric": "a"}, want{http.StatusBadRequest, "", ""}},
		{"test 4", map[string]string{"type": "counter", "metric": "a"}, want{http.StatusOK, "", "text/plain; charset=utf-8"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/value", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			request.SetPathValue("type", test.pathValues["type"])
			if metric, ok := test.pathValues["metric"]; ok {
				request.SetPathValue("metric", metric)
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			})

			handlerToTest := CheckMetricsTypeMw(nextHandler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestCheckPostMethodMw(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string // добавляем название тестов
		method       string
		body         models.Metrics // добавляем тело запроса в табличные тесты
		expectedCode int
		expectedBody string
	}{
		{
			name:         "test 1",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "test 2",
			method:       http.MethodPost,
			body:         models.Metrics{ID: "metric1", MType: "gauge", Delta: nil, Value: nil},
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "test 3",
			method:       http.MethodPost,
			body:         models.Metrics{ID: "metric2", MType: "counter", Delta: nil, Value: nil},
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			request := httptest.NewRequest(tc.method, "/update", bytes.NewReader(body))
			// создаём новый Recorder
			w := httptest.NewRecorder()

			handlerToTest := CheckPostMethodMw(handler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tc.expectedCode, res.StatusCode, "Response code didn't match expected")

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			respBody, err := io.ReadAll(res.Body)
			require.NoError(t, err, "error making HTTP request")

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(respBody))
			}
		})
	}
}

func TestGZIPMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var metrics models.Metrics
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&metrics); err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		ms := &memory.MemStorage{
			Gauges:   make(map[string]models.Gauge),
			Counters: make(map[string]models.Counter),
		}

		_ = ms.AddCounter(metrics.ID, models.Counter(0))
		delta, _ := ms.GetCounter(metrics.ID)
		metrics.Delta = &delta

		rw.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(rw)
		if err := enc.Encode(metrics); err != nil {
			return
		}
	})
	srv := httptest.NewServer(GZIPMiddleware(handler))
	defer srv.Close()

	requestBody := `{
        "id": "PollCount",
        "type": "counter"
    }`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)
		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(r)

		require.NoError(t, err, "error making HTTP request")
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		log.Println(resp.Header.Get("Accept-Encoding"))

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err, "error in new reader")

		_, err = io.ReadAll(zr)
		require.NoError(t, err, "error in read all")
	})
}
