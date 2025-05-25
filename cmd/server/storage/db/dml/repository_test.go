package dml

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateGauge(t *testing.T) {
	tests := []struct {
		name       string
		gaugeName  string
		gaugeValue models.Gauge
		wantErr    bool
	}{
		{
			name:       "success create gauge",
			gaugeName:  "metric1",
			gaugeValue: models.Gauge(1.1),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mock sqlmock.Sqlmock
			DBRepository.Database, mock, _ = sqlmock.New()
			defer DBRepository.Database.Close()

			mock.ExpectExec("^INSERT INTO gauge *").
				WithArgs("metric1", float64(1.1)).
				WillReturnResult(sqlmock.NewResult(1, 1))
			_, err := CreateOrUpdateGauge(&DBRepository, tt.gaugeName, tt.gaugeValue)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success close",
			wantErr: false,
		},
		{
			name:    "error close",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			dbr := &Repository{
				Database: db,
			}

			if tt.wantErr {
				mock.ExpectClose().WillReturnError(fmt.Errorf("close error"))
			} else {
				mock.ExpectClose()
			}

			err = dbr.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_Open(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{
			name:    "invalid dsn",
			dsn:     "invalid://dsn",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.DatabaseDSN = tt.dsn
			rep := &Repository{}
			_, err := rep.Open()
			assert.NoError(t, err)
		})
	}
}

func TestRepository_PingContext(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create repository with mock DB
	repo := &Repository{Database: db}

	t.Run("successful ping", func(t *testing.T) {
		mock.ExpectPing()

		err = repo.PingContext(context.Background())
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestRepository_SetDatabase(t *testing.T) {
	// Setup test cases
	tests := []struct {
		dbr     *Repository
		name    string
		wantErr bool
	}{
		{
			dbr:     &Repository{},
			name:    "invalid connection",
			wantErr: true,
		},
	}
	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.DatabaseDSN = "invalid_dsn"

			err := tt.dbr.SetDatabase()
			assert.NoError(t, err)
			assert.Nil(t, tt.dbr.Database)

			if tt.dbr.Database != nil {
				_ = tt.dbr.Database.Close()
			}
		})
	}
}

func TestRepository_QueryRowContext(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repository{Database: db}

	ctx := context.Background()
	query := "SELECT 1"
	args := []any{}

	row := repo.QueryRowContext(ctx, query, args...)
	assert.NotNil(t, row)
}

func TestRepository_QueryContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repository{Database: db}

	ctx := context.Background()
	query := "SELECT 1"
	args := "metric1"
	rows := sqlmock.NewRows([]string{"value"}).AddRow("1")
	mock.ExpectQuery("^SELECT *").WithArgs("metric1").WillReturnRows(rows)

	row, errQueryContext := repo.QueryContext(ctx, query, args)
	assert.NotNil(t, row)
	assert.NoError(t, errQueryContext)
}
