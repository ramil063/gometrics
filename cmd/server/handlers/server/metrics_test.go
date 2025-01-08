package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/cmd/agent/storage"
)

func TestPrepareMetricsValues(t *testing.T) {
	type args struct {
		ms Storager
		m  storage.Monitor
	}
	a := args{
		ms: NewMemStorage(),
		m:  storage.NewMonitor(),
	}
	tests := []struct {
		name string
		args args
	}{
		{"test 1", a},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PrepareMetricsValues(tt.args.ms, tt.args.m)
			assert.NoError(t, err)
			pc, _ := tt.args.ms.GetCounter("PollCount")
			assert.Equal(t, pc, int64(tt.args.m.PollCount+1))
		})
	}
}

func TestSaveMetricsPerTime(t *testing.T) {
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
			SaveMetricsPerTime(tt.args.workTime, tt.args.ticker, GetStorage("", ""))
		})
	}
}
