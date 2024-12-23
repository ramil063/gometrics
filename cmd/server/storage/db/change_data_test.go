package db

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/models"
)

func TestStorage_AddCounter(t *testing.T) {
	var mock sqlmock.Sqlmock
	dml.DBRepository.Database, mock, _ = sqlmock.New()
	defer dml.DBRepository.Database.Close()

	mock.ExpectExec("^INSERT INTO counter *").
		WithArgs("metric1", int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	type args struct {
		name  string
		value models.Counter
	}

	var a = args{
		name:  "metric1",
		value: models.Counter(1),
	}

	tests := []struct {
		name string
		args args
	}{
		{"test 1", a},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			err := s.AddCounter(tt.args.name, tt.args.value)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_SetGauge(t *testing.T) {
	var mock sqlmock.Sqlmock
	dml.DBRepository.Database, mock, _ = sqlmock.New()
	defer dml.DBRepository.Database.Close()

	mock.ExpectExec("^INSERT INTO gauge *").
		WithArgs("metric1", float64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	type args struct {
		name  string
		value models.Gauge
	}

	var a = args{
		name:  "metric1",
		value: models.Gauge(1),
	}
	tests := []struct {
		name string
		args args
	}{
		{"test 1", a},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{}
			err := s.SetGauge(tt.args.name, tt.args.value)
			assert.NoError(t, err)
		})
	}
}
