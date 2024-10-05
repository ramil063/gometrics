package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ramil063/gometrics/cmd/server/storage"
	"io"
	"net/http"
	"strconv"
)

var ms = storage.NewMemStorage()

// Router маршрутизация
func Router() chi.Router {
	r := chi.NewRouter()

	r.Use(CheckMethodMw)

	r.Route("/", func(r chi.Router) {
		r.Route("/update", func(r chi.Router) {
			r.Use(CheckActionsMw)
			r.Use(CheckUpdateMetricsMw)
			r.Route("/{type}/{metric}/{value}", func(r chi.Router) {
				r.Post("/", update)
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Use(CheckActionsMw)
			r.Use(CheckValueMetricsMw)
			r.Route("/{type}/{metric}", func(r chi.Router) {
				r.Get("/", getValue)
			})
		})
		r.Get("/", home)
	})
	return r
}

// Update метод обновления данных для метрик
func update(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")
	metricValue := r.PathValue("value")

	switch metricType {
	case "gauge":
		value, _ := strconv.ParseFloat(metricValue, 64)
		ms.SetGauge(metricName, storage.Gauge(value))
	case "counter":
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		ms.AddCounter(metricName, storage.Counter(value))
	}
	_, err := io.WriteString(rw, "")
	if err != nil {
		return
	}
}

// getValue метод получения данных из метрики
func getValue(rw http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("metric")

	var err error
	switch metricType {
	case "gauge":
		value, ok := ms.GetGauge(metricName)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err = io.WriteString(rw, strconv.FormatFloat(value, 'f', 1, 64))
	case "counter":
		value, ok := ms.GetCounter(metricName)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err = io.WriteString(rw, strconv.FormatInt(value, 10))
	}
	if err != nil {
		return
	}
}

// Home метод получения данных из всех метрик
func home(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	bodyGauge := ""
	for key, g := range ms.GetGauges() {
		str := strconv.FormatFloat(float64(g), 'f', 1, 64)
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
