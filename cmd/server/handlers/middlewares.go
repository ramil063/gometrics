package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// CheckMethodMw middleware для проверки метода запроса
func CheckMethodMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем только POST, GET запросы с определенным заголовком
		if r.Method != http.MethodPost && r.Method != http.MethodGet {
			log.Println("Error in method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Неверный запрос"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// CheckActionsMw middleware для проверки экшенов
func CheckActionsMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		regUpdate := regexp.MustCompile(`/update/`)
		notUpdateURL := len(regUpdate.Find([]byte(r.URL.Path)))

		regValue := regexp.MustCompile(`/value/`)
		notValueURL := len(regValue.Find([]byte(r.URL.Path)))
		//разрешаем только update, value экшны
		if notUpdateURL == 0 && notValueURL == 0 {
			log.Println("Error in action")
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// CheckUpdateMetricsMw middleware для проверки типа метрик и их значений
func CheckUpdateMetricsMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		re := regexp.MustCompile(`/(gauge|counter)/`)
		issetCorrectType := len(re.Find([]byte(r.URL.Path))) > 1
		// если тип метрики указан неверно
		if !issetCorrectType {
			log.Println("Error in metric type")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// по умолчанию и название метрики и значение неправильное
		issetMetricData := false
		issetCorrectValue := false

		re = regexp.MustCompile(`/update/gauge/`)
		isGaugeUpdate := len(re.Find([]byte(r.URL.Path))) > 0
		// если обновляем метрику gauge
		if isGaugeUpdate {
			metricValue := r.URL.Path[len("/update/gauge/"):]
			metricData := strings.Split(metricValue, "/")
			issetMetricData = len(metricData) == 2
			// если есть и название метрики и значение
			if issetMetricData {
				// если значение верно
				if _, err := strconv.ParseFloat(metricData[1], 64); err == nil {
					issetCorrectValue = true
				}
			}
		}

		re = regexp.MustCompile(`/update/counter/`)
		isCounterUpdate := len(re.Find([]byte(r.URL.Path))) > 0
		// если обновляем метрику counter
		if isCounterUpdate {
			metricValue := r.URL.Path[len("/update/counter/"):]
			metricData := strings.Split(metricValue, "/")
			issetMetricData = len(metricData) == 2
			// если есть и название метрики и значение
			if issetMetricData {
				// если значение верно
				if _, err := strconv.ParseInt(metricData[1], 10, 64); err == nil {
					issetCorrectValue = true
				}
			}
		}

		if !issetMetricData {
			log.Println("Error in metric data")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if !issetCorrectValue {
			log.Println("Error in metric value")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			next.ServeHTTP(w, r)
		}
	})
}

// CheckValueMetricsMw middleware для проверки типа метрик
func CheckValueMetricsMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// по умолчанию и название метрики и значение неправильное
		issetMetricName := false

		re := regexp.MustCompile(`/value/gauge/`)
		isGetGauge := len(re.Find([]byte(r.URL.Path))) > 0
		// если обновляем метрику gauge
		if isGetGauge {
			metricName := r.URL.Path[len("/value/gauge/"):]
			// если есть название метрики
			if metricName != "" {
				issetMetricName = true
			}
		}

		re = regexp.MustCompile(`/value/counter/`)
		isCounterUpdate := len(re.Find([]byte(r.URL.Path))) > 0
		// если обновляем метрику counter
		if isCounterUpdate {
			metricName := r.URL.Path[len("/value/counter/"):]
			// если есть название метрики
			if metricName != "" {
				issetMetricName = true
			}
		}

		if !issetMetricName {
			log.Println("Error in metric name")
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			next.ServeHTTP(w, r)
		}
	})
}
