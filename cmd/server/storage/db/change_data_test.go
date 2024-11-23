package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/models"
	"testing"
)

func TestStorage_AddCounter(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dml.DBRepository.Database, _, _ = sqlmock.New()
	defer dml.DBRepository.Database.Close()

	type fields struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		name  string
		value models.Counter
	}
	var f = fields{
		Gauges:   map[string]models.Gauge{},
		Counters: map[string]models.Counter{},
	}

	var a = args{
		name:  "metric1",
		value: models.Counter(1),
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"test 1", f, a},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			s.AddCounter(tt.args.name, tt.args.value)
		})
	}
}

func TestStorage_SetGauge(t *testing.T) {

	dml.DBRepository.Database, _, _ = sqlmock.New()
	defer dml.DBRepository.Database.Close()

	type fields struct {
		Gauges   map[string]models.Gauge
		Counters map[string]models.Counter
	}
	type args struct {
		name  string
		value models.Gauge
	}
	var f = fields{
		Gauges:   map[string]models.Gauge{},
		Counters: map[string]models.Counter{},
	}

	var a = args{
		name:  "metric1",
		value: models.Gauge(1),
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"test 1", f, a},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			s.SetGauge(tt.args.name, tt.args.value)
		})
	}
}
