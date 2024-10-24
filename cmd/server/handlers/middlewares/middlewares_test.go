package middlewares

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
