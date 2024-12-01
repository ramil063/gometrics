package db

import (
	"context"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"time"

	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

type Storage struct {
	Gauges   map[string]models.Gauge
	Counters map[string]models.Counter
}

func Init(dbr dml.DataBaser) error {
	var err error

	if err = CheckPing(dbr); err != nil {
		logger.WriteErrorLog("Database ping error", err.Error())
		return err
	}

	err = CreateTables(dbr)
	return err
}

func CheckPing(dbr dml.DataBaser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return dbr.PingContext(ctx)
}

func CreateTables(dbr dml.DataBaser) error {
	var err error

	createTablesSQL := `
	
	CREATE TABLE IF NOT EXISTS public.gauge
	(
	    id    serial constraint gauge_pk primary key,
	    name  varchar          not null constraint gauge_pk_2 unique,
	    value double precision not null
	);
	comment on table public.gauge is 'Gauge метрики';
	comment on column public.gauge.name is 'Название метрики';
	comment on column public.gauge.value is 'Значение метрики';
	        
	CREATE TABLE IF NOT EXISTS public.counter
	(
	    id    serial constraint counter_pk primary key,
	    name  varchar not null constraint counter_pk_2 unique,
	    value bigint not null
	);
	comment on table public.counter is 'Counter метрики';
	comment on column public.counter.name is 'Название метрики';
	comment on column public.counter.value is 'Значение метрики';`

	_, err = dbr.ExecContext(context.Background(), createTablesSQL)
	return err
}
