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

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/cmd/server/storage/memory"
	"github.com/ramil063/gometrics/internal/hash"
	"github.com/ramil063/gometrics/internal/models"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

func TestCheckMethodMw(t *testing.T) {
	type want struct {
		contentType string
		code        int
	}
	tests := []struct {
		name        string
		method      string
		contentType string
		want        want
	}{
		{
			name:        "test 1",
			method:      http.MethodGet,
			contentType: "text-plain",
			want: want{
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
		{
			name:        "test 2",
			method:      http.MethodPost,
			contentType: "application/json",
			want: want{
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
		{
			name:        "test 3",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
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
		response    string
		contentType string
		code        int
	}
	tests := []struct {
		name       string
		pathValues map[string]string
		want       want
	}{
		{
			name:       "test 1",
			pathValues: map[string]string{"type": "gauge", "metric": "a", "value": "1"},
			want: want{
				response:    "",
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
		{
			name:       "test 2",
			pathValues: map[string]string{"type": "gauge", "metric": "a"},
			want: want{
				response:    "",
				contentType: "",
				code:        http.StatusBadRequest,
			},
		},
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
		response    string
		contentType string
		code        int
	}
	tests := []struct {
		name       string
		pathValues map[string]string
		want       want
	}{
		{
			name:       "test 1",
			pathValues: map[string]string{"type": "gauge", "metric": "a"},
			want: want{
				response:    "",
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
		{
			name:       "test 2",
			pathValues: map[string]string{"type": "gauge1", "metric": "a"},
			want: want{
				response:    "",
				contentType: "",
				code:        http.StatusBadRequest,
			},
		},
		{
			name:       "test 3",
			pathValues: map[string]string{"type": "counter"},
			want: want{
				response:    "",
				contentType: "",
				code:        http.StatusNotFound,
			},
		},
		{
			name:       "test 4",
			pathValues: map[string]string{"type": "counter", "metric": "a"},
			want: want{
				response:    "",
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
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
		response    string
		contentType string
		code        int
	}
	tests := []struct {
		name       string
		pathValues map[string]string
		want       want
	}{
		{
			name:       "test 1",
			pathValues: map[string]string{"type": "gauge", "metric": "a"},
			want: want{
				response:    "",
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
		{
			name:       "test 4",
			pathValues: map[string]string{"type": "counter"},
			want: want{
				response:    "",
				contentType: "",
				code:        http.StatusBadRequest,
			},
		},
		{
			name:       "test 4",
			pathValues: map[string]string{"type": "counter", "metric": "a"},
			want: want{
				response:    "",
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
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
		response    string
		contentType string
		code        int
	}
	tests := []struct {
		name       string
		pathValues map[string]string
		want       want
	}{
		{
			name:       "test 1",
			pathValues: map[string]string{"type": "gauge", "metric": "a"},
			want: want{
				response:    "",
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
		{
			name:       "test 2",
			pathValues: map[string]string{"type": "counter1", "metric": "a"},
			want: want{
				response:    "",
				contentType: "",
				code:        http.StatusBadRequest,
			},
		},
		{
			name:       "test 3",
			pathValues: map[string]string{"type": "counter", "metric": "a"},
			want: want{
				response:    "",
				contentType: "text/plain; charset=utf-8",
				code:        http.StatusOK,
			},
		},
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
		body         models.Metrics // добавляем тело запроса в табличные тесты
		name         string         // добавляем название тестов
		method       string
		expectedBody string
		expectedCode int
	}{
		{
			expectedBody: "",
			name:         "test 1",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			body:         models.Metrics{ID: "metric1", MType: "gauge", Delta: nil, Value: nil},
			name:         "test 2",
			method:       http.MethodPost,
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
		{
			body:         models.Metrics{ID: "metric2", MType: "counter", Delta: nil, Value: nil},
			name:         "test 3",
			method:       http.MethodPost,
			expectedBody: "",
			expectedCode: http.StatusOK,
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

func TestCheckHashMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(handler)
	defer srv.Close()

	handlers.HashKey = "test"

	testCases := []struct {
		body         models.Metrics // добавляем тело запроса в табличные тесты
		name         string         // добавляем название тестов
		method       string
		expectedBody string
		expectedCode int
	}{
		{
			body:         models.Metrics{ID: "metric1", MType: "gauge", Delta: nil, Value: nil},
			name:         "test 1",
			method:       http.MethodPost,
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
		{
			body:         models.Metrics{ID: "metric2", MType: "counter", Delta: nil, Value: nil},
			name:         "test 2",
			method:       http.MethodPost,
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			request := httptest.NewRequest(tc.method, "/update", bytes.NewReader(body))

			request.Header.Set("HashSHA256", hash.CreateSha256(body, handlers.HashKey))
			request.Header.Set("Content-Type", "application/json")
			// создаём новый Recorder
			w := httptest.NewRecorder()

			handlerToTest := CheckHashMiddleware(handler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tc.expectedCode, res.StatusCode, "Response code didn't match expected")

			defer res.Body.Close()
		})
	}
}

// MockDecryptor реализует интерфейс Decryptor для тестов
type MockDecryptor struct {
	decryptFunc func([]byte) ([]byte, error)
}

func (m *MockDecryptor) Decrypt(data []byte) ([]byte, error) {
	return m.decryptFunc(data)
}

func TestDecryptMiddleware(t *testing.T) {
	// Сохраняем оригинальный decryptor и восстанавливаем после теста
	originalDecryptor := crypto.DefaultDecryptor
	defer func() {
		crypto.DefaultDecryptor = originalDecryptor
	}()

	tests := []struct {
		name              string
		expectedBody      string
		decryptor         crypto.Decryptor
		requestBody       []byte
		expectedStatus    int
		expectNextHandler bool
	}{
		{
			name:              "No decryptor - pass through",
			expectedBody:      "original",
			decryptor:         nil,
			requestBody:       []byte("original"),
			expectedStatus:    http.StatusOK,
			expectNextHandler: true,
		},
		{
			name:         "Successful decryption",
			expectedBody: "decrypted",
			decryptor: &MockDecryptor{
				decryptFunc: func(data []byte) ([]byte, error) {
					return []byte("decrypted"), nil
				},
			},
			requestBody:       []byte("encrypted"),
			expectedStatus:    http.StatusOK,
			expectNextHandler: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crypto.DefaultDecryptor = tt.decryptor

			body := bytes.NewReader(tt.requestBody)

			req, err := http.NewRequest("POST", "/", body)
			if err != nil {
				t.Fatal(err)
			}

			nextHandlerCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextHandlerCalled = true
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				w.Write(body)
			})

			rr := httptest.NewRecorder()
			middleware := DecryptMiddleware(nextHandler)
			middleware.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectNextHandler != nextHandlerCalled {
				t.Errorf("next handler called: got %v want %v",
					nextHandlerCalled, tt.expectNextHandler)
			}

			if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
				t.Errorf("unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

// TestCheckTrustedIP tests the CheckTrustedIP middleware for various trust scenarios
func TestCheckTrustedIP(t *testing.T) {
	// Save and restore the original TrustedSubnet after test
	originalTrustedSubnet := handlers.TrustedSubnet
	defer func() { handlers.TrustedSubnet = originalTrustedSubnet }()

	tests := []struct {
		name           string
		trustedSubnet  string
		realIP         string
		expectedStatus int
	}{
		{
			name:           "No TrustedSubnet set, should pass",
			trustedSubnet:  "",
			realIP:         "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "TrustedSubnet set, X-Real-IP in subnet, should pass",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "192.168.1.42",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "TrustedSubnet set, X-Real-IP not in subnet, should fail",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "10.0.0.1",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "TrustedSubnet set, missing X-Real-IP header, should fail",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.TrustedSubnet = tt.trustedSubnet

			req := httptest.NewRequest("GET", "/", nil)
			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}

			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			})

			rr := httptest.NewRecorder()
			CheckTrustedIP(next).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.True(t, called, "next handler should be called for allowed requests")
			} else {
				assert.False(t, called, "next handler should not be called for forbidden requests")
			}
		})
	}
}
