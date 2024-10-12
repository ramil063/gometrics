package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_update(t *testing.T) {
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
				update(rw, req, ms)
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
	ms := NewMemStorage()
	ts := httptest.NewServer(Router(ms))
	defer ts.Close()

	type want struct {
		code     int
		response string
	}
	testsG := []struct {
		name string
		url  string
		want want
	}{
		{"gauge test 1", "/value/gauge/a", want{200, ""}},
		{"gauge test 2", "/value/gauge/a1", want{404, ""}},
		{"gauge test 3", "/value/gauge/", want{404, ""}},
		{"gauge test 4", "/value/gauge/testUnknown80", want{404, ""}},
	}
	for _, test := range testsG {
		t.Run(test.name, func(t *testing.T) {
			ms.SetGauge("a", 1.1)
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
		{"counter test 1", "/value/counter/a", want{200, ""}},
		{"counter test 2", "/value/counter/a1", want{404, ""}},
		{"counter test 3", "/value/counter/", want{404, ""}},
	}
	for _, test := range testsC {
		t.Run(test.name, func(t *testing.T) {
			ms.AddCounter("a", 1)
			resp, _ := testRequest(t, ts, "GET", test.url)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
		})
	}
}

func Test_home(t *testing.T) {
	ms := NewMemStorage()
	ts := httptest.NewServer(Router(ms))
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{"test 1", "/", want{200, "", "text/html; charset=utf-8"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms.SetGauge("a", 1.1)
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
		{"test 1", "*storage.MemStorage"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStorage()
			assert.Equalf(t, tt.want, reflect.ValueOf(ms).Type().String(), "NewMemStorage()")
		})
	}
}
