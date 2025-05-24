package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/internal/models"
)

func TestNewMonitor(t *testing.T) {
	tests := []struct {
		want *Monitor
		name string
	}{
		{
			name: "check monitor",
			want: &Monitor{PollCount: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMonitor()
			assert.Equal(t, m.PollCount, tt.want.PollCount)
		})
	}
}

func TestMonitor_InitCPUutilizationValue(t *testing.T) {
	type fields struct {
		CPUutilization map[int]models.Gauge
	}
	tests := []struct {
		fields fields
		name   string
	}{
		{
			name: "test 1",
			fields: fields{
				CPUutilization: map[int]models.Gauge{1: models.Gauge(1.1)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Monitor{
				CPUutilization: tt.fields.CPUutilization,
			}
			m.InitCPUutilizationValue()
			assert.NotEqual(t, nil, m.CPUutilization)
			assert.NotEqual(t, tt.fields.CPUutilization[1], m.CPUutilization[1])
		})
	}
}

func TestMonitor_StoreCPUutilizationValue(t *testing.T) {
	type fields struct {
		CPUutilization map[int]models.Gauge
	}
	type args struct {
		key   int
		value models.Gauge
	}
	tests := []struct {
		fields fields
		name   string
		args   args
	}{
		{
			name:   "test 1",
			fields: fields{CPUutilization: map[int]models.Gauge{1: models.Gauge(1.1)}},
			args:   args{key: 1, value: models.Gauge(2.2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Monitor{
				CPUutilization: tt.fields.CPUutilization,
			}
			m.StoreCPUutilizationValue(tt.args.key, tt.args.value)
			assert.Equal(t, tt.args.value, m.CPUutilization[tt.args.key])
		})
	}
}

func TestMonitor_GetAllCPUutilization(t *testing.T) {
	type fields struct {
		CPUutilization map[int]models.Gauge
	}
	tests := []struct {
		want   map[int]models.Gauge
		fields fields
		name   string
	}{
		{
			name:   "test 1",
			want:   map[int]models.Gauge{1: models.Gauge(1.1)},
			fields: fields{CPUutilization: map[int]models.Gauge{1: models.Gauge(1.1)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Monitor{
				CPUutilization: tt.fields.CPUutilization,
			}
			assert.Equalf(t, tt.want, m.GetAllCPUutilization(), "GetAllCPUutilization()")
		})
	}
}
