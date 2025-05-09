package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/go-chi/chi/v5"

	agentStorage "github.com/ramil063/gometrics/cmd/agent/storage"
	"github.com/ramil063/gometrics/cmd/server/handlers/middlewares"
	"github.com/ramil063/gometrics/cmd/server/storage/db"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/models"
)

var MaxSaverWorkTime = 900000

type Gauger interface {
	SetGauge(name string, value models.Gauge) error
	GetGauge(name string) (float64, error)
	GetGauges() (map[string]models.Gauge, error)
}

type Counterer interface {
	AddCounter(name string, value models.Counter) error
	GetCounter(name string) (int64, error)
	GetCounters() (map[string]models.Counter, error)
}

type Storager interface {
	Gauger
	Counterer
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

	r.Get("/ping", ping)

	r.Route("/updates", func(r chi.Router) {
		r.Use(middlewares.CheckHashMiddleware)
		updatesHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
			updates(rw, r, s)
		}
		r.With(middlewares.CheckPostMethodMw).Post("/", updatesHandlerFunction)
	})

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

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// Для heap/goroutine/block:
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))

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
		err := ms.SetGauge(metricName, models.Gauge(value))
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Gauge")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "counter":
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		err := ms.AddCounter(metricName, models.Counter(value))
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Counter")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
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

	switch metricType {
	case "gauge":
		value, err := ms.GetGauge(metricName)
		if err != nil {
			log.Println("Error value of gauge is not Ok")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain")
		_, err = io.WriteString(rw, strconv.FormatFloat(value, 'f', -1, 64))
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Gauge")
		}
	case "counter":
		value, err := ms.GetCounter(metricName)
		if err != nil {
			log.Println("Error value of counter is not Ok")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain")
		_, err = io.WriteString(rw, strconv.FormatInt(value, 10))
		if err != nil {
			logger.WriteErrorLog(err.Error(), "Counter")
		}
	}
}

// Home метод получения данных из всех метрик
func home(rw http.ResponseWriter, r *http.Request, ms Storager) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "text/html")

	bodyGauge := ""
	gauges, err := ms.GetGauges()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "Gauge")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	for key, g := range gauges {
		str := strconv.FormatFloat(float64(g), 'f', -1, 64)
		bodyGauge += `<div>` + key + `-` + str + `</div>`
	}
	bodyCounters := ""
	counters, err := ms.GetCounters()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "GetCounters")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	for key, c := range counters {
		str := strconv.FormatInt(int64(c), 10)
		bodyCounters += `<div>` + key + `-` + str + `</div>`
	}

	_, err = io.WriteString(rw, `
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
		logger.WriteErrorLog(err.Error(), "WriteString")
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

	logMsg, _ := json.Marshal(metrics)
	logger.WriteInfoLog("request body in update/", string(logMsg))

	switch metrics.MType {
	case "gauge":
		err := s.SetGauge(metrics.ID, models.Gauge(*metrics.Value))
		if err != nil {
			logger.WriteErrorLog(err.Error(), "SetGauge")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "counter":
		err := s.AddCounter(metrics.ID, models.Counter(*metrics.Delta))
		if err != nil {
			logger.WriteErrorLog(err.Error(), "AddCounter")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		newCounter, err := s.GetCounter(metrics.ID)
		if err != nil {
			logger.WriteDebugLog(err.Error(), "GetCounter")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		metrics.Delta = &newCounter
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

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

	switch metrics.MType {
	case "gauge":
		value, err := s.GetGauge(metrics.ID)
		if err != nil {
			logger.WriteInfoLog(err.Error(), "GetGauge ID:"+metrics.ID)
			err = s.SetGauge(metrics.ID, 0)
			if err != nil {
				logger.WriteInfoLog(err.Error(), "SetGauge ID:"+metrics.ID)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			value = 0
		}
		metrics.Value = &value
	case "counter":
		err := s.AddCounter(metrics.ID, models.Counter(0))
		if err != nil {
			logger.WriteInfoLog(err.Error(), "AddCounter ID:"+metrics.ID)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		delta, err := s.GetCounter(metrics.ID)
		if err != nil {
			logger.WriteInfoLog(err.Error(), "GetCounter ID:"+metrics.ID)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		metrics.Delta = &delta
	}

	m := agentStorage.NewMonitor()
	err := PrepareMetricsValues(s, m)

	if err != nil {
		logger.WriteErrorLog(err.Error(), "PrepareMetricsValues getValueMetricsJSON")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	if err := enc.Encode(metrics); err != nil {
		logger.WriteErrorLog("error encoding response", err.Error())
		return
	}
	logger.WriteDebugLog("", "sending HTTP 200 response")
}

// Home метод получения данных из всех метрик
func ping(rw http.ResponseWriter, r *http.Request) {
	rep, err := dml.NewRepository()
	if err != nil {
		logger.WriteErrorLog("Database storage ping error", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = db.CheckPing(rep); err != nil {
		logger.WriteErrorLog("Database storage ping error", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// Home метод получения данных из всех метрик
func updates(rw http.ResponseWriter, r *http.Request, dbs Storager) {
	var metrics []models.Metrics
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&metrics)

	if err != nil {
		logger.WriteDebugLog("cannot decode request JSON body", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	//only for autotests
	//logMsg, _ := json.Marshal(metrics)
	//logger.WriteInfoLog("request body in updates/", string(logMsg))

	result, err := updateMetrics(dbs, metrics)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	if err = enc.Encode(result); err != nil {
		logger.WriteErrorLog("error encoding response", err.Error())
		return
	}

	logger.WriteDebugLog("", "sending HTTP 200 response")
}

func updateMetrics(dbs Storager, metrics []models.Metrics) ([]models.Metrics, error) {
	result := make([]models.Metrics, 0, len(metrics))

	for _, m := range metrics {
		current := m

		switch current.MType {
		case "gauge":
			if current.Value == nil {
				zero := 0.0
				current.Value = &zero
			}
			if err := dbs.SetGauge(current.ID, models.Gauge(*current.Value)); err != nil {
				logger.WriteErrorLog(err.Error(), "SetGauge ID:"+current.ID)
				return nil, err
			}
		case "counter":
			if current.Delta == nil {
				zero := int64(0)
				current.Delta = &zero
			}
			if err := dbs.AddCounter(current.ID, models.Counter(*current.Delta)); err != nil {
				logger.WriteErrorLog(err.Error(), "AddCounter ID:"+current.ID)
				return nil, err
			}
			newCounter, err := dbs.GetCounter(current.ID)
			if err != nil {
				logger.WriteErrorLog(err.Error(), "GetCounter ID:"+m.ID)
				return nil, err
			}
			current.Delta = &newCounter
		}

		result = append(result, current)
	}

	return result, nil
}
