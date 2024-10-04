package storage

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type MemStorageMock struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
}

func TestMemStorage_AddCounter(t *testing.T) {
	type args struct {
		name  string
		value Counter
	}
	tests := []struct {
		name           string
		MemStorageMock MemStorageMock
		args           args
		want           args
	}{
		{
			"test 1",
			MemStorageMock{
				Counters: map[string]Counter{"counter1": 1},
			},
			args{
				name:  "counter1",
				value: 1,
			},
			args{
				name:  "counter1",
				value: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Counters: tt.MemStorageMock.Counters,
			}
			ms.AddCounter(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.value, ms.Counters[tt.args.name])
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	type args struct {
		name  string
		value Gauge
	}
	tests := []struct {
		name           string
		MemStorageMock MemStorageMock
		args           args
		want           args
	}{
		{
			"test 1",
			MemStorageMock{
				Gauges: map[string]Gauge{"gauge1": 1.1},
			},
			args{
				name:  "gauge1",
				value: 3.5,
			},
			args{
				name:  "gauge1",
				value: 3.5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Gauges: tt.MemStorageMock.Gauges,
			}
			ms.SetGauge(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.value, ms.Gauges[tt.args.name])
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
