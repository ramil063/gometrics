package middlewares

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/ramil063/gometrics/internal/logger"
)

// CheckMethodMw middleware для проверки метода запроса
func CheckMethodMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем только POST, GET запросы
		if r.Method != http.MethodPost && r.Method != http.MethodGet {
			logger.Log.Debug("Incorrect method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Incorrect method"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CheckPostMethodMw middleware для проверки метода запроса
func CheckPostMethodMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем только POST запросы
		if r.Method != http.MethodPost {
			logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CheckUpdateMetricsNameMw middleware для проверки имени метрик
func CheckUpdateMetricsNameMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("metric") == "" {
			logger.Log.Debug("Error metric name is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CheckMetricsTypeMw middleware для проверки типа метрик
func CheckMetricsTypeMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("type") != "gauge" && r.PathValue("type") != "counter" {
			logger.Log.Debug("Error in metric type (allowed 'gauge' or 'counter') got " + r.PathValue("type"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CheckUpdateMetricsValueMw middleware для проверки значения метрик
func CheckUpdateMetricsValueMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("value") == "" {
			logger.Log.Debug("Error in metric value")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// по умолчанию и название метрики и значение неправильное
		issetMetricData := false
		issetCorrectValue := false

		// если обновляем метрику gauge
		if r.PathValue("type") == "gauge" {
			// если есть и название метрики и значение
			if r.PathValue("metric") != "" && r.PathValue("value") != "" {
				issetMetricData = true
				// если значение верно
				if _, err := strconv.ParseFloat(r.PathValue("value"), 64); err == nil {
					issetCorrectValue = true
				}
			}
		}

		// если обновляем метрику counter
		if r.PathValue("type") == "counter" {
			// если есть и название метрики и значение
			if r.PathValue("metric") != "" && r.PathValue("value") != "" {
				issetMetricData = true
				// если значение верно
				if _, err := strconv.ParseInt(r.PathValue("value"), 10, 64); err == nil {
					issetCorrectValue = true
				}
			}
		}

		if !issetMetricData {
			logger.Log.Debug("Error in metric data(update)")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if !issetCorrectValue {
			logger.Log.Debug("Error in metric value(update)")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// CheckValueMetricsMw middleware для проверки типа метрик
func CheckValueMetricsMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("type") != "gauge" && r.PathValue("type") != "counter" {
			logger.Log.Debug("Error in metric type")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.PathValue("metric") == "" {
			logger.Log.Debug("Error in metric name")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
