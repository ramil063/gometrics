package storage

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/cmd/agent/storage"
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
	filename := "../../../internal/storage/files/test.json"
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
			err := SaveMonitorPerSeconds(tt.args.workTime, tt.args.ticker, 1, filename)
			assert.NoError(t, err)
		})
	}
}

func TestSaveMonitor(t *testing.T) {
	filename := "../../../internal/storage/files/test.json"

	tests := []struct {
		name     string
		filename string
	}{
		{"test 1", filename},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, SaveMonitor(tt.filename))
		})
	}
}
