package db

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/models"
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

			got, err := s.GetCounter(tt.args.name)
			if got != tt.want {
				t.Errorf("GetCounter() got = %v, want %v", got, tt.want)
			}
			if err != nil {
				t.Errorf("GetCounter() got1 = %v, want %v", err, nil)
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
	mock.ExpectQuery("^SELECT name, value FROM counter").WithArgs().WillReturnRows(rows)

	tests := []struct {
		want map[string]models.Counter
		name string
	}{
		{
			want: map[string]models.Counter{"metric1": 1, "metric2": 2},
			name: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			got, err := s.GetCounters()
			assert.NoError(t, err)
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
			if got1 != nil {
				t.Errorf("GetGauge() got1 = %v, want %v", got1, nil)
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
		want map[string]models.Gauge
		name string
	}{
		{
			want: map[string]models.Gauge{"metric1": 1.1, "metric2": 2.2},
			name: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			got, err := s.GetGauges()
			assert.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGauges() = %v, want %v", got, tt.want)
			}
		})
	}
}
