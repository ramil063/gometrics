package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/ramil063/gometrics/internal/logger"
)

type DataBaser interface {
	Init(ps string) error
	CheckPing() error
}

type DB struct {
	Ptr *sql.DB
}

var Database DB

func (db *DB) Init(ps string) error {
	database, err := sql.Open("pgx", ps)
	if err != nil {
		logger.WriteErrorLog("DB open error", err.Error())
		return err
	}
	db.Ptr = database

	if err = db.CheckPing(); err != nil {
		logger.WriteErrorLog("DB ping error", err.Error())
		return err
	}

	Database.Ptr = database
	return nil
}

func (db *DB) CheckPing() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.Ptr.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
