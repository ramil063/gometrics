package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	agentStorage "github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/cmd/server/handlers/middlewares"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

var MaxSaverWorkTime = 900000

type Gauger interface {
	SetGauge(name string, value models.Gauge)
	GetGauge(name string) (float64, bool)
	GetGauges() map[string]models.Gauge
}

type Counterer interface {
	AddCounter(name string, value models.Counter)
	GetCounter(name string) (int64, bool)
	GetCounters() map[string]models.Counter
}

type Storager interface {
	Gauger
	Counterer
}

type DBStorager interface {
	Storager
	CreateOrUpdateCounter(name string, value models.Counter) (sql.Result, error)
	CreateOrUpdateGauge(name string, value models.Gauge) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Open() (*sql.DB, error)
	Close() error
	PingContext(ctx context.Context) error
	SetDatabase() error
}

// Router маршрутизация
func Router(s Storager) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.ResponseLogger)
	r.Use(logger.RequestLogger)
	r.Use(middlewares.GZIPMiddleware)
	r.Use(middlewares.CheckMethodMw)

	homeHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
		home(rw, r, s)
	}
	r.Get("/", homeHandlerFunction)

	pingHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
		ping(rw, r, s)
	}
	r.Get("/ping", pingHandlerFunction)

	r.Route("/update", func(r chi.Router) {
		r.Route("/{type}/{metric}", func(r chi.Router) {
			r.Use(middlewares.CheckMetricsTypeMw)
			r.Use(middlewares.CheckUpdateMetricsNameMw)
			updateHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
				update(rw, req, s)
			}
			r.With(middlewares.CheckUpdateMetricsValueMw).Post("/", updateHandlerFunction)
			r.With(middlewares.CheckUpdateMetricsValueMw).Post("/{value}", updateHandlerFunction)
		})

		updateMetricsJSONHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
			updateMetricsJSON(rw, req, s)
		}
		r.With(middlewares.CheckPostMethodMw).Post("/", updateMetricsJSONHandlerFunction)
	})
	r.Route("/value", func(r chi.Router) {
		r.Route("/{type}/{metric}", func(r chi.Router) {
			r.Use(middlewares.CheckMetricsTypeMw)
			r.Use(middlewares.CheckValueMetricsMw)
			getValueHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
				getValue(rw, req, s)
			}
			r.Get("/", getValueHandlerFunction)
		})

		getValueMetricsJSONHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
			getValueMetricsJSON(rw, req, s)
		}
		r.With(middlewares.CheckPostMethodMw).Post("/", getValueMetricsJSONHandlerFunction)
	})
	return r
}

// Update метод обновления данных для метрик
func update(rw http.ResponseWriter, r *http.Request, ms Storager) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")
	metricValue := r.PathValue("value")

	switch metricType {
	case "gauge":
		value, _ := strconv.ParseFloat(metricValue, 64)
		ms.SetGauge(metricName, models.Gauge(value))
	case "counter":
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		ms.AddCounter(metricName, models.Counter(value))
	}
	_, err := io.WriteString(rw, "")
	if err != nil {
		return
	}
}

// getValue метод получения данных из метрики
func getValue(rw http.ResponseWriter, r *http.Request, ms Storager) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")

	var err error
	switch metricType {
	case "gauge":
		value, ok := ms.GetGauge(metricName)
		if !ok {
			log.Println("Error value of gauge is not Ok")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain")
		_, err = io.WriteString(rw, strconv.FormatFloat(value, 'f', -1, 64))
	case "counter":
		value, ok := ms.GetCounter(metricName)
		if !ok {
			log.Println("Error value of counter is not Ok")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain")
		_, err = io.WriteString(rw, strconv.FormatInt(value, 10))
	}
	if err != nil {
		return
	}
}

// Home метод получения данных из всех метрик
func home(rw http.ResponseWriter, r *http.Request, ms Storager) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "text/html")

	bodyGauge := ""
	for key, g := range ms.GetGauges() {
		str := strconv.FormatFloat(float64(g), 'f', -1, 64)
		bodyGauge += `<div>` + key + `-` + str + `</div>`
	}
	bodyCounters := ""
	for key, c := range ms.GetCounters() {
		str := strconv.FormatInt(int64(c), 10)
		bodyCounters += `<div>` + key + `-` + str + `</div>`
	}

	_, err := io.WriteString(rw, `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Все метрики</title>
</head>
	<body>
		<t2>Gauge</t2>
		`+bodyGauge+`
		<t2>Counters</t2>
		`+bodyCounters+`
	</body>
</html>
`)
	if err != nil {
		return
	}
}

// updateMetricsJSON метод обновления данных для метрик через json
func updateMetricsJSON(rw http.ResponseWriter, r *http.Request, s Storager) {

	// десериализуем запрос в структуру модели
	logger.WriteDebugLog("", "decoding request")
	var metrics models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.WriteDebugLog("cannot decode request JSON body", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	logMsg, _ := json.Marshal(metrics)
	logger.WriteInfoLog("request body in update/", string(logMsg))

	switch metrics.MType {
	case "gauge":
		s.SetGauge(metrics.ID, models.Gauge(*metrics.Value))
	case "counter":
		s.AddCounter(metrics.ID, models.Counter(*metrics.Delta))
		newCounter, _ := s.GetCounter(metrics.ID)
		metrics.Delta = &newCounter
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(metrics); err != nil {
		logger.WriteErrorLog("error encoding response", err.Error())
		return
	}

	logger.WriteDebugLog("", "sending HTTP 200 response")
}

// getValueMetricsJSON метод обновления данных для метрик через json
func getValueMetricsJSON(rw http.ResponseWriter, r *http.Request, s Storager) {
	// десериализуем запрос в структуру модели
	logger.WriteDebugLog("message", "decoding request")
	var metrics models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.WriteDebugLog("cannot decode request JSON body", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.WriteInfoLog("request body in value/", "metrics{ID, MType}"+metrics.ID+","+metrics.MType)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	switch metrics.MType {
	case "gauge":
		value, ok := s.GetGauge(metrics.ID)
		if !ok {
			logger.WriteInfoLog("Error value of gauge is not Ok", "ID:"+metrics.ID)
			s.SetGauge(metrics.ID, 0)
			value = 0
		}
		metrics.Value = &value
	case "counter":
		s.AddCounter(metrics.ID, models.Counter(0))
		delta, ok := s.GetCounter(metrics.ID)
		if !ok {
			logger.WriteInfoLog("Error value of counter is not Ok", "ID:"+metrics.ID)
		}
		metrics.Delta = &delta
	}

	m := agentStorage.NewMonitor()
	PrepareMetricsValues(s, m)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(metrics); err != nil {
		logger.WriteErrorLog("error encoding response", err.Error())
		return
	}
	logger.WriteDebugLog("", "sending HTTP 200 response")
}

// Home метод получения данных из всех метрик
func ping(rw http.ResponseWriter, r *http.Request, ms Storager) {
	if err := db.CheckPing(dml.NewRepository()); err != nil {
		logger.WriteErrorLog("Database storage ping error", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
