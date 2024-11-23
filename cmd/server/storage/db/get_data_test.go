package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/models"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestStorage_GetCounter(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dml.DBRepository.Database = db
	rows := sqlmock.NewRows([]string{"value"}).AddRow("1")
	mock.ExpectQuery("^SELECT value FROM counter WHERE name = *").WithArgs("metric1").WillReturnRows(rows)

	type args struct {
		name string
	}
	var a = args{
		name: "metric1",
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 bool
	}{
		{"test 1", a, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}

			got, got1 := s.GetCounter(tt.args.name)
			if got != tt.want {
				t.Errorf("GetCounter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCounter() got1 = %v, want %v", got1, tt.want1)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all expectations were met: %v", err)
			}
		})
	}
}

func TestStorage_GetCounters(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dml.DBRepository.Database = db

	rows := sqlmock.NewRows([]string{"name", "value"}).AddRow("metric1", "1").AddRow("metric2", "2")
	mock.ExpectQuery("^SELECT name, value FROM counter").WillReturnRows(rows)

	tests := []struct {
		name string
		want map[string]models.Counter
	}{
		{"test 1", map[string]models.Counter{"metric1": 1, "metric2": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			got := s.GetCounters()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStorage_GetGauge(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dml.DBRepository.Database = db
	rows := sqlmock.NewRows([]string{"value"}).AddRow("1.1")
	mock.ExpectQuery("^SELECT value FROM gauge WHERE name = *").WithArgs("metric1").WillReturnRows(rows)

	type args struct {
		name string
	}
	tests := []struct {
		name   string
		args   args
		want   float64
		wantOk bool
	}{
		{"test 1", args{"metric1"}, 1.1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			got, got1 := s.GetGauge(tt.args.name)
			if got != tt.want {
				t.Errorf("GetGauge() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("GetGauge() got1 = %v, want %v", got1, tt.wantOk)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all expectations were met: %v", err)
			}
		})
	}
}

func TestStorage_GetGauges(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dml.DBRepository.Database = db

	rows := sqlmock.NewRows([]string{"name", "value"}).AddRow("metric1", "1.1").AddRow("metric2", "2.2")
	mock.ExpectQuery("^SELECT name, value FROM gauge").WillReturnRows(rows)

	tests := []struct {
		name string
		want map[string]models.Gauge
	}{
		{"test 1", map[string]models.Gauge{"metric1": 1.1, "metric2": 2.2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			if got := s.GetGauges(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGauges() = %v, want %v", got, tt.want)
			}
		})
	}
}
