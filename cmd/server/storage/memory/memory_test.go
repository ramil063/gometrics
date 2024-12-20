package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/internal/models"
)

type MemStorageMock struct {
	Gauges   map[string]models.Gauge
	Counters map[string]models.Counter
}

func TestMemStorage_AddCounter(t *testing.T) {
	type args struct {
		name  string
		value models.Counter
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
				Counters: map[string]models.Counter{"counter1": 1},
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
			err := ms.AddCounter(tt.args.name, tt.args.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.value, ms.Counters[tt.args.name])
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	type args struct {
		name  string
		value models.Gauge
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
				Gauges: map[string]models.Gauge{"gauge1": 1.1},
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
			err := ms.SetGauge(tt.args.name, tt.args.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.value, ms.Gauges[tt.args.name])
		})
	}
}
