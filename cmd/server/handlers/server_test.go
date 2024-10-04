package handlers

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"test 1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewServer(":45566")
			go func() {
				time.Sleep(2 * time.Second)
				err := srv.Shutdown(context.Background())
				if err != nil {
					return
				}
			}()
			err := srv.Run()
			if ((err != nil) && !errors.Is(err, http.ErrServerClosed)) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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

			nextHandler := http.HandlerFunc(update)

			handlerToTest := CheckMetricsMw(nextHandler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.statusCode, res.StatusCode)
			// получаем и проверяем тело запроса
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}

func TestNewServer(t *testing.T) {
	type args struct {
		adr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test 1", args{adr: "127.0.0.1:45566"}, `handlers.Server`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := reflect.ValueOf(NewServer(tt.args.adr)).Type().String()
			assert.Equal(t, tt.want, actual, "NewServer(%v)", tt.args.adr)
		})
	}
}
