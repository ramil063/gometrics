package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// CheckMethodMw middleware для проверки метода запроса
func CheckMethodMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем только POST-запросы с определенным заголовком
		if r.Method != http.MethodPost {
			log.Println("Error in method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Неверный запрос"))
		} else if r.Header.Get("Content-Type") != "text/plain" {
			log.Println("Error in context-type")
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
		re := regexp.MustCompile(`/update/`)
		notUpdateURL := len(re.Find([]byte(r.URL.Path)))
		//разрешаем только update экшн
		if notUpdateURL == 0 {
			log.Println("Error in action")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CheckMetricsMw middleware для проверки типа метрик и их значений
func CheckMetricsMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		re := regexp.MustCompile(`/(gauge|counter)/`)
		issetCorrectType := len(re.Find([]byte(r.URL.Path))) > 1
		// если тип метрики указане неверно
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
