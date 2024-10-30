package storage

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/cmd/server/handlers"
)

func TestGetMonitor(t *testing.T) {
	m := storage.Monitor{}

	tests := []struct {
		name string
		want storage.Monitor
	}{
		{"test 1", m},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := GetMonitor(false)
			assert.Equalf(t, reflect.ValueOf(tt.want).Kind(), reflect.ValueOf(mc).Kind(), "GetMonitor()")
		})
	}
}

func TestSaveMonitorPerSeconds(t *testing.T) {
	type args struct {
		workTime int
		ticker   *time.Ticker
	}
	a := args{
		workTime: 1,
		ticker:   time.NewTicker(1 * time.Second),
	}
	tests := []struct {
		name string
		args args
	}{
		{"test 1", a},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveMonitorPerSeconds(tt.args.workTime, tt.args.ticker, 1, "../../../"+handlers.FileStoragePath)
			assert.NoError(t, err)
		})
	}
}
