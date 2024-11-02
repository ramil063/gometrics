package file

import (
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFStorage_AddCounter(t *testing.T) {
	var f = FStorage{
		Gauges:   map[string]models.Gauge{},
		Counters: map[string]models.Counter{"met1": 1},
	}

	handlers.FileStoragePath = "../../../../internal/storage/files/test.json"
	_, err := NewReader(handlers.FileStoragePath)

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
		value models.Counter
	}
	tests := []struct {
		name   string
		fields FStorage
		args   args
	}{
		{"test 1", f, args{name: "met1", value: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			s.AddCounter(tt.args.name, tt.args.value)
			got, ok := s.GetCounter(tt.args.name)
			assert.Equal(t, ok, true)
			assert.Equal(t, int64(tt.fields.Counters[tt.args.name]+tt.args.value), got)
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
		name   string
		fields FStorage
		args   args
		want   int64
		want1  bool
	}{
		{"test 1", f, args{name: "met1"}, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			got, got1 := s.GetCounter(tt.args.name)
			assert.Equalf(t, tt.want, got, "GetCounter(%v)", tt.args.name)
			assert.Equalf(t, tt.want1, got1, "GetCounter(%v)", tt.args.name)
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
		name   string
		fields FStorage
		want   map[string]models.Counter
	}{
		{"test 1", f, map[string]models.Counter{"met1": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			assert.Equalf(t, tt.want, s.GetCounters(), "GetCounters()")
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
		name   string
		fields FStorage
		args   args
		want   float64
		want1  bool
	}{
		{"test 1", f, args{name: "met1"}, 1.1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			got, got1 := s.GetGauge(tt.args.name)
			assert.Equalf(t, tt.want, got, "GetGauge(%v)", tt.args.name)
			assert.Equalf(t, tt.want1, got1, "GetGauge(%v)", tt.args.name)
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
		name   string
		fields FStorage
		want   map[string]models.Gauge
	}{
		{"test 1", f, map[string]models.Gauge{"met1": 1.1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			assert.Equalf(t, tt.want, s.GetGauges(), "GetGauges()")
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
		name   string
		fields FStorage
		args   args
		want   map[string]models.Gauge
	}{
		{"test 1", f, args{name: "met1", value: 1.1}, map[string]models.Gauge{"met1": 1.1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			s.SetGauge(tt.args.name, tt.args.value)
			assert.Equalf(t, tt.want, s.GetGauges(), "GetGauges()")
		})
	}
}
