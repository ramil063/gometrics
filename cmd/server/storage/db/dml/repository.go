package dml

import (
	"context"
	"database/sql"
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

type Repository struct {
	Database *sql.DB
}

type DataBaser interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Open() (*sql.DB, error)
	Close() error
	PingContext(ctx context.Context) error
	SetDatabase() error
}

var DBRepository Repository

func (dbr *Repository) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return dbr.Database.ExecContext(ctx, query, args...)
}

func (dbr *Repository) Close() error {
	return dbr.Database.Close()
}

func (dbr *Repository) PingContext(ctx context.Context) error {
	return dbr.Database.PingContext(ctx)
}

func (dbr *Repository) SetDatabase() error {
	database, err := dbr.Open()
	if err != nil {
		logger.WriteErrorLog("Database open error", err.Error())
	}
	dbr.Database = database
	return nil
}

func (dbr *Repository) Open() (*sql.DB, error) {
	return sql.Open("pgx", handlers.DatabaseDSN)
}

func NewRepository() *Repository {
	rep := &Repository{}
	_ = rep.SetDatabase()
	return rep
}

func CreateOrUpdateCounter(dbr *Repository, name string, value models.Counter) (sql.Result, error) {
	var result sql.Result
	row := dbr.Database.QueryRowContext(context.Background(), "SELECT name FROM counter WHERE name = $1", name)

	if row.Err() != nil {
		logger.WriteErrorLog("AddCounter database query error", row.Err().Error())
		return result, row.Err()
	}
	var selectedName string
	_ = row.Scan(&selectedName)

	if selectedName != "" {
		return dbr.ExecContext(
			context.Background(),
			"UPDATE counter SET value = $1 + value WHERE name = $2",
			int64(value),
			name)
	}
	return dbr.ExecContext(
		context.Background(),
		"INSERT INTO counter (name, value) VALUES ($1, $2)",
		name,
		float64(value))
}

func CreateOrUpdateGauge(dbr *Repository, name string, value models.Gauge) (sql.Result, error) {
	var result sql.Result

	row := dbr.Database.QueryRowContext(context.Background(), "SELECT name FROM gauge WHERE name = $1", name)

	if row.Err() != nil {
		logger.WriteErrorLog("SetGauge database query error", row.Err().Error())
		return result, row.Err()
	}
	var selectedName string
	_ = row.Scan(&selectedName)

	if selectedName != "" {
		return dbr.ExecContext(
			context.Background(),
			"UPDATE gauge SET value = $1 WHERE name = $2",
			float64(value),
			name)
	}
	return dbr.ExecContext(
		context.Background(),
		"INSERT INTO gauge (name, value) VALUES ($1, $2)",
		name,
		float64(value))
}
