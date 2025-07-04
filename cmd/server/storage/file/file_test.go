package file

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

func TestFStorage_AddCounter(t *testing.T) {
	handlers.FileStoragePath = "../../../../internal/storage/files/test.json"

	type storage struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		name  string
		value models.Counter
	}
	tests := []struct {
		name string
		storage
		args args
	}{
		{
			"test 1",
			storage{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{"met1": 1},
			},
			args{
				name:  "met1",
				value: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.storage.Gauges,
				Counters: tt.storage.Counters,
				mx:       sync.RWMutex{},
			}
			err := s.AddCounter(tt.args.name, tt.args.value)
			assert.NoError(t, err)
			got, err := s.GetCounter(tt.args.name)
			assert.NoError(t, err)
			assert.Equal(t, got+int64(tt.args.value), got+int64(tt.storage.Counters[tt.args.name]))
		})
	}
}

func TestFStorage_GetCounter(t *testing.T) {
	var f = FStorage{
		Gauges:   map[string]models.Gauge{},
		Counters: map[string]models.Counter{"met1": 1},
	}

	handlers.FileStoragePath = "../../../../internal/storage/files/test.json"
	Writer, err := NewWriter(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error create metrics writer", err.Error())
	}
	defer Writer.Close()

	err = Writer.WriteMetrics(&f)
	if err != nil {
		logger.WriteErrorLog("error write metrics", err.Error())
	}
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"test 1", args{name: "met1"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   f.Gauges,
				Counters: f.Counters,
				mx:       sync.RWMutex{},
			}
			got, err := s.GetCounter(tt.args.name)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetCounter(%v)", tt.args.name)
		})
	}
}

func TestFStorage_GetCounters(t *testing.T) {
	var f = FStorage{
		Gauges:   map[string]models.Gauge{},
		Counters: map[string]models.Counter{"met1": 1},
	}

	handlers.FileStoragePath = "../../../../internal/storage/files/test.json"
	Writer, err := NewWriter(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error create metrics writer", err.Error())
	}
	defer Writer.Close()

	err = Writer.WriteMetrics(&f)
	if err != nil {
		logger.WriteErrorLog("error write metrics", err.Error())
	}
	tests := []struct {
		want map[string]models.Counter
		name string
	}{
		{
			want: map[string]models.Counter{"met1": 1},
			name: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   f.Gauges,
				Counters: f.Counters,
				mx:       sync.RWMutex{},
			}
			got, err := s.GetCounters()
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetCounters()")
		})
	}
}

func TestFStorage_GetGauge(t *testing.T) {
	var f = FStorage{
		Gauges:   map[string]models.Gauge{"met1": 1.1},
		Counters: map[string]models.Counter{},
	}

	handlers.FileStoragePath = "../../../../internal/storage/files/test.json"
	Writer, err := NewWriter(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error create metrics writer", err.Error())
	}
	defer Writer.Close()

	err = Writer.WriteMetrics(&f)
	if err != nil {
		logger.WriteErrorLog("error write metrics", err.Error())
	}
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"test 1", args{name: "met1"}, 1.1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   f.Gauges,
				Counters: f.Counters,
				mx:       sync.RWMutex{},
			}
			got, err := s.GetGauge(tt.args.name)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetGauge(%v)", tt.args.name)
		})
	}
}

func TestFStorage_GetGauges(t *testing.T) {
	var f = FStorage{
		Gauges:   map[string]models.Gauge{"met1": 1.1},
		Counters: map[string]models.Counter{},
	}

	handlers.FileStoragePath = "../../../../internal/storage/files/test.json"
	Writer, err := NewWriter(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error create metrics writer", err.Error())
	}
	defer Writer.Close()

	err = Writer.WriteMetrics(&f)
	if err != nil {
		logger.WriteErrorLog("error write metrics", err.Error())
	}
	tests := []struct {
		want map[string]models.Gauge
		name string
	}{
		{
			want: map[string]models.Gauge{"met1": 1.1},
			name: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   f.Gauges,
				Counters: f.Counters,
				mx:       sync.RWMutex{},
			}
			got, err := s.GetGauges()
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetGauges()")
		})
	}
}

func TestFStorage_SetGauge(t *testing.T) {
	var f = FStorage{
		Gauges:   map[string]models.Gauge{"met1": 1.1},
		Counters: map[string]models.Counter{},
	}

	handlers.FileStoragePath = "../../../../internal/storage/files/test.json"
	Writer, err := NewWriter(handlers.FileStoragePath)
	if err != nil {
		logger.WriteErrorLog("error create metrics writer", err.Error())
	}
	defer Writer.Close()

	err = Writer.WriteMetrics(&f)
	if err != nil {
		logger.WriteErrorLog("error write metrics", err.Error())
	}
	type args struct {
		name  string
		value models.Gauge
	}
	tests := []struct {
		want map[string]models.Gauge
		name string
		args args
	}{
		{
			want: map[string]models.Gauge{"met1": 1.1},
			name: "test 1",
			args: args{
				name:  "met1",
				value: 1.1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   f.Gauges,
				Counters: f.Counters,
			}
			err = s.SetGauge(tt.args.name, tt.args.value)
			assert.NoError(t, err)
			got, err := s.GetGauges()
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetGauges()")
		})
	}
}

func TestFStorage_StoreGaugeValue(t *testing.T) {
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
				Gauges:   map[string]models.Gauge{"met1": 1.1},
				Counters: map[string]models.Counter{},
			},
			args{
				key:   "met1",
				value: 1.1,
			},
		},
	}
	for _, tt := range tests {
		s := &FStorage{
			Gauges:   tt.s.Gauges,
			Counters: tt.s.Counters,
			mx:       sync.RWMutex{},
		}
		s.StoreGaugeValue(tt.args.key, tt.args.value)
		got, err := s.GetGaugeValue(tt.args.key)
		assert.NoError(t, err)
		assert.Equalf(t, tt.args.value, got, "GetGaugeValue")
	}
}

func TestFStorage_GetGaugeValue(t *testing.T) {
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
			s := &FStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreGaugeValue(tt.args.key, tt.want)
			got, err := s.GetGaugeValue(tt.args.key)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetGaugeValue(%v)", tt.args.key)
		})
	}
}

func TestFStorage_GetAllGauges(t *testing.T) {
	type storage struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		key   string
		value models.Gauge
	}
	tests := []struct {
		want map[string]models.Gauge
		s    storage
		name string
		args args
	}{
		{
			want: map[string]models.Gauge{"met1": 1.1},
			s: storage{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			name: "test 1",
			args: args{
				key:   "met1",
				value: 1.1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreGaugeValue(tt.args.key, tt.args.value)
			assert.Equalf(t, tt.want, s.GetAllGauges(), "GetAllGauges()")
		})
	}
}

func TestFStorage_StoreCounterValue(t *testing.T) {
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
			s := &FStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreCounterValue(tt.args.key, tt.args.value)
			got, err := s.GetCounterValue(tt.args.key)
			assert.NoError(t, err)
			assert.Equalf(t, tt.args.value, got, "GetCounterValue")
		})
	}
}

func TestFStorage_GetCounterValue(t *testing.T) {
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
		want models.Counter
	}{
		{
			"test 1",
			storage{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			args{
				key: "met1",
			},
			models.Counter(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreCounterValue(tt.args.key, tt.want)
			got, err := s.GetCounterValue(tt.args.key)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "GetCounterValue(%v)", tt.args.key)
		})
	}
}

func TestFStorage_GetAllCounters(t *testing.T) {
	type storage struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		key   string
		value models.Counter
	}
	tests := []struct {
		want map[string]models.Counter
		s    storage
		name string
		args args
	}{
		{
			want: map[string]models.Counter{"met1": 1},
			s: storage{
				Gauges:   map[string]models.Gauge{},
				Counters: map[string]models.Counter{},
			},
			name: "test 1",
			args: args{
				key:   "met1",
				value: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.s.Gauges,
				Counters: tt.s.Counters,
				mx:       sync.RWMutex{},
			}
			s.StoreCounterValue(tt.args.key, tt.args.value)
			assert.Equalf(t, tt.want, s.GetAllCounters(), "GetAllCounters()")
		})
	}
}
