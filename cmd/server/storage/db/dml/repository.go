package dml

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	internalErrors "github.com/ramil063/gometrics/internal/errors"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

type Repository struct {
	Database *sql.DB
}

type DataBaser interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Open() (*sql.DB, error)
	Close() error
	PingContext(ctx context.Context) error
	SetDatabase() error
}

var DBRepository Repository

func (dbr *Repository) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	result, err := dbr.Database.ExecContext(ctx, query, args...)
	if err != nil {
		var pgconnErr *pgconn.PgError
		if errors.As(err, &pgconnErr) && pgerrcode.IsConnectionException(pgconnErr.Code) {
			result, err = retryExecContext(dbr, internalErrors.TriesTimes, ctx, query, args)
			if err == nil {
				return result, nil
			}
		}
		return nil, internalErrors.NewDBError(err)
	}
	return result, nil
}

func (dbr *Repository) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	row := dbr.Database.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		var pgconnErr *pgconn.PgError
		if errors.As(row.Err(), &pgconnErr) && pgerrcode.IsConnectionException(pgconnErr.Code) {
			row = retryQueryRowContext(dbr, internalErrors.TriesTimes, ctx, query, args)
		}
	}
	return row
}

func (dbr *Repository) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := dbr.Database.QueryContext(ctx, query, args...)

	if err != nil {
		var pgconnErr *pgconn.PgError
		if errors.As(err, &pgconnErr) && pgerrcode.IsConnectionException(pgconnErr.Code) {
			rows, err = retryQueryContext(dbr, internalErrors.TriesTimes, ctx, query, args)
			if err != nil {
				return nil, internalErrors.NewDBError(err)
			}
		}
	}
	return rows, nil
}

func (dbr *Repository) Close() error {
	err := dbr.Database.Close()
	if err != nil {
		return internalErrors.NewDBError(err)
	}
	return nil
}

func (dbr *Repository) PingContext(ctx context.Context) error {
	err := dbr.Database.PingContext(ctx)
	if err != nil {
		var pgconnErr *pgconn.PgError
		if errors.As(err, &pgconnErr) && pgerrcode.IsConnectionException(pgconnErr.Code) {
			err = retryPing(dbr, ctx, internalErrors.TriesTimes)
			if err != nil {
				return internalErrors.NewDBError(err)
			}
		}
	}
	return err
}

func (dbr *Repository) SetDatabase() error {
	database, err := dbr.Open()
	if err != nil {
		logger.WriteErrorLog("Database open error", err.Error())
		return internalErrors.NewDBError(err)
	}
	dbr.Database = database
	return nil
}

func (dbr *Repository) Open() (*sql.DB, error) {
	result, err := sql.Open("pgx", handlers.DatabaseDSN)
	if err != nil {
		var pgconnErr *pgconn.PgError
		if errors.As(err, &pgconnErr) && pgerrcode.IsConnectionException(pgconnErr.Code) {
			result, err = retryOpen("pgx", handlers.DatabaseDSN, internalErrors.TriesTimes)
			if err != nil {
				return nil, internalErrors.NewDBError(err)
			}
		}
	}
	return result, nil
}

func NewRepository() (*Repository, error) {
	rep := &Repository{}
	err := rep.SetDatabase()
	return rep, err
}

func CreateOrUpdateCounter(dbr *Repository, name string, value models.Counter) (sql.Result, error) {
	exec, err := dbr.ExecContext(
		context.Background(),
		"INSERT INTO counter (name, value) VALUES ($1, $2) "+
			"ON CONFLICT (name) "+
			"DO UPDATE SET value = $2 + counter.value "+
			"WHERE counter.name = $1",
		name,
		int64(value))

	if err != nil {
		return nil, internalErrors.NewDBError(err)
	}
	return exec, nil
}

func CreateOrUpdateGauge(dbr *Repository, name string, value models.Gauge) (sql.Result, error) {
	exec, err := dbr.ExecContext(
		context.Background(),
		"INSERT INTO gauge (name, value) VALUES ($1, $2) "+
			"ON CONFLICT (name) "+
			"DO UPDATE SET value = $2 "+
			"WHERE gauge.name = $1",
		name,
		float64(value))
	if err != nil {
		return nil, internalErrors.NewDBError(err)
	}
	return exec, nil
}

func retryQueryRowContext(dbr *Repository, tries []int, ctx context.Context, query string, args ...any) *sql.Row {
	var row *sql.Row
	for try := 0; try < len(tries); try++ {
		time.Sleep(time.Duration(tries[try]) * time.Second)
		row = dbr.Database.QueryRowContext(ctx, query, args)
		if row.Err() == nil {
			break
		}
	}
	return row
}

func retryQueryContext(dbr *Repository, tries []int, ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	for try := 0; try < len(tries); try++ {
		time.Sleep(time.Duration(tries[try]) * time.Second)
		rows, err = dbr.Database.QueryContext(ctx, query, args)
		if err == nil {
			break
		}
	}
	return rows, err
}

func retryExecContext(dbr *Repository, tries []int, ctx context.Context, query string, args ...any) (sql.Result, error) {
	var result sql.Result
	var err error
	for try := 0; try < len(tries); try++ {
		time.Sleep(time.Duration(tries[try]) * time.Second)
		result, err = dbr.Database.ExecContext(ctx, query, args...)
		if err == nil {
			break
		}
	}
	return result, err
}

func retryOpen(driverName, dataSourceName string, tries []int) (*sql.DB, error) {
	var result *sql.DB
	var err error
	for try := 0; try < len(tries); try++ {
		time.Sleep(time.Duration(tries[try]) * time.Second)
		result, err = sql.Open(driverName, dataSourceName)
		if err == nil {
			break
		}
	}
	return result, err
}

func retryPing(dbr *Repository, ctx context.Context, tries []int) error {
	var err error
	for try := 0; try < len(tries); try++ {
		time.Sleep(time.Duration(tries[try]) * time.Second)
		err = dbr.Database.PingContext(ctx)
		if err == nil {
			break
		}
	}
	return err
}
