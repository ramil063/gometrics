package handlers

import (
	"github.com/ramil063/gometrics/cmd/server/storage"
	"net/http"
	"strconv"
)

var ms = storage.NewMemStorage()

type Server struct {
	*http.Server
}

// NewServer метод для запуска сервера
func NewServer(adr string) Server {
	return Server{
		&http.Server{
			Addr:    adr,
			Handler: nil,
		},
	}
}

// Run метод для запуска сервера
func (s Server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/update/{type}/{metric}/{value}", Conveyor(http.HandlerFunc(update), CheckMethodMw, CheckMetricsMw, CheckActionsMw))
	s.Handler = mux
	return s.ListenAndServe()
}

// update метод обновления данных для метрик
func update(w http.ResponseWriter, r *http.Request) {
	metricName := r.PathValue("metric")
	metricValue := r.PathValue("value")

	switch metricName {
	case "gauge":
		value, _ := strconv.ParseFloat(metricValue, 64)
		ms.SetGauge(metricName, storage.Gauge(value))
	case "counter":
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		ms.AddCounter(metricName, storage.Counter(value))
	}
}
