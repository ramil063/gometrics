package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/models"
)

func TestStorage_AddCounter(t *testing.T) {
	var mock sqlmock.Sqlmock
	dml.DBRepository.Database, mock, _ = sqlmock.New()
	defer dml.DBRepository.Database.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("1")
	mock.ExpectQuery("^SELECT name FROM counter WHERE name = *").WithArgs("metric1").WillReturnRows(rows)

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
	var mock sqlmock.Sqlmock
	dml.DBRepository.Database, mock, _ = sqlmock.New()
	defer dml.DBRepository.Database.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("1")
	mock.ExpectQuery("^SELECT name FROM gauge WHERE name = *").WithArgs("metric1").WillReturnRows(rows)

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
