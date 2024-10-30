package server

import (
	"encoding/json"
	"fmt"
	"github.com/ramil063/gometrics/cmd/server/handlers"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	agentStorage "github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/cmd/server/handlers/middlewares"
	"github.com/ramil063/gometrics/cmd/server/storage"
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

func NewMemStorage() Storager {
	return &storage.MemStorage{
		Gauges:   make(map[string]models.Gauge),
		Counters: make(map[string]models.Counter),
	}
}

// Router маршрутизация
func Router(ms Storager) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.ResponseLogger)
	r.Use(logger.RequestLogger)
	r.Use(middlewares.GZIPMiddleware)
	r.Use(middlewares.CheckMethodMw)
	r.Use(middlewares.SaveMonitorToFile)

	homeHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
		home(rw, r, ms)
	}
	r.Get("/", homeHandlerFunction)

	r.Route("/update", func(r chi.Router) {
		r.Route("/{type}/{metric}", func(r chi.Router) {
			r.Use(middlewares.CheckMetricsTypeMw)
			r.Use(middlewares.CheckUpdateMetricsNameMw)
			updateHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
				update(rw, req, ms)
			}
			r.With(middlewares.CheckUpdateMetricsValueMw).Post("/", updateHandlerFunction)
			r.With(middlewares.CheckUpdateMetricsValueMw).Post("/{value}", updateHandlerFunction)
		})

		updateMetricsJSONHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
			updateMetricsJSON(rw, req, ms)
		}
		r.With(middlewares.CheckPostMethodMw).Post("/", updateMetricsJSONHandlerFunction)
	})
	r.Route("/value", func(r chi.Router) {
		r.Route("/{type}/{metric}", func(r chi.Router) {
			r.Use(middlewares.CheckMetricsTypeMw)
			r.Use(middlewares.CheckValueMetricsMw)
			getValueHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
				getValue(rw, req, ms)
			}
			r.Get("/", getValueHandlerFunction)
		})

		getValueMetricsJSONHandlerFunction := func(rw http.ResponseWriter, req *http.Request) {
			getValueMetricsJSON(rw, req, ms)
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
func updateMetricsJSON(rw http.ResponseWriter, r *http.Request, ms Storager) {

	// десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var metrics models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	logMsg, _ := json.Marshal(metrics)
	logger.Log.Info("", zap.String("request body in update/", string(logMsg)))

	switch metrics.MType {
	case "gauge":
		ms.SetGauge(metrics.ID, models.Gauge(*metrics.Value))
	case "counter":
		ms.AddCounter(metrics.ID, models.Counter(*metrics.Delta))
		newCounter, _ := ms.GetCounter(metrics.ID)
		metrics.Delta = &newCounter
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(metrics); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}

	logger.Log.Debug("sending HTTP 200 response")
}

// getValueMetricsJSON метод обновления данных для метрик через json
func getValueMetricsJSON(rw http.ResponseWriter, r *http.Request, ms Storager) {
	// десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var metrics models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Log.Info("request body in value/", zap.String("metrics{ID, MType}", metrics.ID+","+metrics.MType))

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	switch metrics.MType {
	case "gauge":
		value, ok := ms.GetGauge(metrics.ID)
		if !ok {
			logger.Log.Info("Error value of gauge is not Ok ID:" + metrics.ID)
			ms.SetGauge(metrics.ID, 0)
			value = 0
		}
		metrics.Value = &value
	case "counter":
		ms.AddCounter(metrics.ID, models.Counter(0))
		delta, ok := ms.GetCounter(metrics.ID)
		if !ok {
			logger.Log.Info("Error value of counter is not Ok ID:" + metrics.ID)
		}
		metrics.Delta = &delta
	}

	m := storage.GetMonitor(handlers.Restore)
	PrepareStorageValues(ms, m)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(metrics); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}

func PrepareStorageValues(ms Storager, m agentStorage.Monitor) {
	v := reflect.ValueOf(m)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		metricID := typeOfS.Field(i).Name
		metricValue, _ := strconv.ParseFloat(fmt.Sprintf("%v", v.Field(i).Interface()), 64)

		if typeOfS.Field(i).Name == "PollCount" {
			ms.AddCounter(metricID, models.Counter(1))
		} else {
			ms.SetGauge(metricID, models.Gauge(metricValue))
		}
	}
}
