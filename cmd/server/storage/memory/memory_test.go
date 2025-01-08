package memory

import (
	"sync"
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

func TestMemStorage_GetAllCounters(t *testing.T) {
	type storage struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		key   string
		value models.Counter
	}
	tests := []struct {
		name string
		s    storage
		args args
		want map[string]models.Counter
	}{
		{
			"test 1",
			storage{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			args{
				key:   "met1",
				value: 1,
			},
			map[string]models.Counter{"met1": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreCounterValue(tt.args.key, tt.args.value)
			assert.Equalf(t, tt.want, s.GetAllCounters(), "GetAllCounters()")
		})
	}
}

func TestMemStorage_GetAllGauges(t *testing.T) {
	type storage struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		key   string
		value models.Gauge
	}
	tests := []struct {
		name string
		s    storage
		args args
		want map[string]models.Gauge
	}{
		{
			"test 1",
			storage{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			args{
				key:   "met1",
				value: 1.1,
			},
			map[string]models.Gauge{"met1": 1.1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreGaugeValue(tt.args.key, tt.args.value)
			assert.Equalf(t, tt.want, s.GetAllGauges(), "GetAllGauges()")
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		name  string
		value models.Counter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			"test 1",
			fields{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			args{
				name:  "met1",
				value: 1,
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
				mx:       sync.RWMutex{},
			}
			err := ms.AddCounter(tt.args.name, tt.args.value)
			assert.NoError(t, err)
			got, err := ms.GetCounter(tt.args.name)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetCounter(%v)", tt.args.name)
		})
	}
}

func TestMemStorage_GetCounterValue(t *testing.T) {
	type storage struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		key   string
		value models.Counter
	}
	tests := []struct {
		name string
		s    storage
		args args
		want models.Counter
	}{
		{
			"test 1",
			storage{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			args{
				key:   "met1",
				value: models.Counter(1),
			},
			models.Counter(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreCounterValue(tt.args.key, tt.args.value)
			got, err := s.GetCounterValue(tt.args.key)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetCounterValue(%v)", tt.args.key)
		})
	}
}

func TestMemStorage_GetCounters(t *testing.T) {
	type fields struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		key   string
		value models.Counter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]models.Counter
	}{
		{
			"test 1",
			fields{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			args{
				key:   "met1",
				value: models.Counter(1),
			},
			map[string]models.Counter{"met1": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
				mx:       sync.RWMutex{},
			}
			err := ms.AddCounter(tt.args.key, tt.args.value)
			assert.NoError(t, err)
			got, err := ms.GetCounters()
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetCounters()")
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		name  string
		value models.Gauge
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			"test 1",
			fields{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			args{
				name:  "met1",
				value: 1.1,
			},
			1.1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
				mx:       sync.RWMutex{},
			}
			err := ms.SetGauge(tt.args.name, tt.args.value)
			assert.NoError(t, err)
			got, err := ms.GetGauge(tt.args.name)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetGauge(%v)", tt.args.name)
		})
	}
}

func TestMemStorage_GetGaugeValue(t *testing.T) {
	type storage struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		key string
	}
	tests := []struct {
		name string
		s    storage
		args args
		want models.Gauge
	}{
		{
			"test 1",
			storage{
				Gauges:   map[string]models.Gauge{"met1": 1.1},
				Counters: map[string]models.Counter{},
			},
			args{
				key: "met1",
			},
			models.Gauge(1.1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			got, err := s.GetGaugeValue(tt.args.key)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetGaugeValue(%v)", tt.args.key)
		})
	}
}

func TestMemStorage_GetGauges(t *testing.T) {
	type fields struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	tests := []struct {
		name string
		f    fields
		want map[string]models.Gauge
	}{
		{
			"test 1",
			fields{
				Gauges:   map[string]models.Gauge{"met1": 1.1},
				Counters: map[string]models.Counter{},
			},
			map[string]models.Gauge{"met1": 1.1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Gauges:   tt.f.Gauges,
				Counters: tt.f.Counters,
				mx:       sync.RWMutex{},
			}
			got, err := ms.GetGauges()
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetGauges()")
		})
	}
}

func TestMemStorage_StoreCounterValue(t *testing.T) {
	type storage struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		key   string
		value models.Counter
	}
	tests := []struct {
		name string
		s    storage
		args args
	}{
		{
			"test 1",
			storage{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			args{
				key:   "met1",
				value: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreCounterValue(tt.args.key, tt.args.value)
			got, err := s.GetCounterValue(tt.args.key)
			assert.NoError(t, err)
			assert.Equalf(t, tt.args.value, got, "StoreCounterValue")
		})
	}
}

func TestMemStorage_StoreGaugeValue(t *testing.T) {
	type storage struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		key   string
		value models.Gauge
	}
	tests := []struct {
		name string
		s    storage
		args args
	}{
		{
			"test 1",
			storage{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			args{
				key:   "met1",
				value: 1.1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreGaugeValue(tt.args.key, tt.args.value)
			got, err := s.GetGaugeValue(tt.args.key)
			assert.NoError(t, err)
			assert.Equalf(t, tt.args.value, got, "StoreGaugeValue")
		})
	}
}
