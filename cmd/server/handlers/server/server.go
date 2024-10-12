package server

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/ramil063/gometrics/cmd/server/handlers/middlewares"
	"github.com/ramil063/gometrics/cmd/server/models"
	"github.com/ramil063/gometrics/cmd/server/storage"
)

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

type MemStorager interface {
	Storager
}

func NewMemStorage() Storager {
	return &storage.MemStorage{
		Gauges:   make(map[string]models.Gauge),
		Counters: make(map[string]models.Counter),
	}
}

// Router маршрутизация
func Router(ms MemStorager) chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.CheckMethodMw)
	homeHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
		home(rw, r, ms)
	}
	r.Get("/", homeHandlerFunction)

	r.Route("/update/{type}", func(r chi.Router) {
		r.Use(middlewares.CheckMetricsTypeMw)
		r.Route("/{metric}", func(r chi.Router) {
			r.Use(middlewares.CheckUpdateMetricsNameMw)
			updateHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
				update(rw, req, ms)
			}
			r.With(middlewares.CheckUpdateMetricsValueMw).Post("/", updateHandlerFunction)
			r.With(middlewares.CheckUpdateMetricsValueMw).Post("/{value}", updateHandlerFunction)
		})
	})
	r.Route("/value/{type}", func(r chi.Router) {
		r.Use(middlewares.CheckMetricsTypeMw)
		r.Route("/{metric}", func(r chi.Router) {
			r.Use(middlewares.CheckValueMetricsMw)
			getValueHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
				getValue(rw, req, ms)
			}
			r.Get("/", getValueHandlerFunction)
		})

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
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err = io.WriteString(rw, strconv.FormatFloat(value, 'f', -1, 64))
	case "counter":
		value, ok := ms.GetCounter(metricName)
		if !ok {
			log.Println("Error value of counter is not Ok")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err = io.WriteString(rw, strconv.FormatInt(value, 10))
	}
	if err != nil {
		return
	}
}

// Home метод получения данных из всех метрик
func home(rw http.ResponseWriter, r *http.Request, ms Storager) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

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
