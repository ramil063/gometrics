package main

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type gauge float64
type counter int64

type Storage interface {
	SetGauge(name string, value gauge)
	GetGauge(name string) (gauge, error)
	AddCounter(name string, value counter)
	GetCounter(name string) (counter, error)
}

type MemStorage struct {
	Gauges   map[string]gauge
	Counters map[string]counter
}

var ms = MemStorage{
	Gauges:   make(map[string]gauge),
	Counters: make(map[string]counter),
}

func (ms *MemStorage) SetGauge(name string, value gauge) {
	ms.Gauges[name] = value
}

func (ms MemStorage) GetGauge(name string) (float64, error) {
	value, ok := ms.Gauges[name]
	if !ok {
		return 0, errors.New("key of gauge is not exist")
	}
	return float64(value), nil
}

func (ms *MemStorage) AddCounter(name string, value counter) {
	oldValue, ok := ms.Counters[name]
	if !ok {
		oldValue = 0
	}
	ms.Counters[name] = oldValue + value
}

func (ms MemStorage) GetCounter(name string) (int64, error) {
	value, ok := ms.Counters[name]
	if !ok {
		return 0, errors.New("key of counter is not exist")
	}
	return int64(value), nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/", updateGauge)
	mux.HandleFunc("/update/counter/", updateCounter)

	middleware3 := checkMetricsMw(mux)
	middleware2 := checkActionsMw(middleware3)
	middleware1 := checkMethodMw(middleware2)

	return http.ListenAndServe(`:8080`, middleware1)
}

// middlware для проверки метода запроса
func checkMethodMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// разрешаем только POST-запросы с определенным заголовком
		if (r.Method == http.MethodPost) &&
			(r.Header.Get("Content-Type") == "text/plain") {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Неверный запрос"))
		}
	})
}

// middlware для проверки экшенов
func checkActionsMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/update/`)
		notUpdateUrl := len(re.Find([]byte(r.URL.Path)))
		//разрешаем только update экшн
		if notUpdateUrl == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// middlware для проверки типа метрик и их значений
func checkMetricsMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		re := regexp.MustCompile(`/(gauge|counter)/`)
		issetCorrectType := len(re.Find([]byte(r.URL.Path))) > 1
		// если тип метрики указане неверно
		if !issetCorrectType {
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
			w.WriteHeader(http.StatusNotFound)
			return
		} else if !issetCorrectValue {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			// установим правильный заголовок для типа данных
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			next.ServeHTTP(w, r)
		}
	})
}

func updateGauge(w http.ResponseWriter, r *http.Request) {
	metricValue := r.URL.Path[len("/update/guage/"):]
	metricData := strings.Split(metricValue, "/")
	metricName := metricData[0]
	value, _ := strconv.ParseFloat(metricData[1], 64)
	ms.SetGauge(metricName, gauge(value))
}

func updateCounter(w http.ResponseWriter, r *http.Request) {
	metricValue := r.URL.Path[len("/update/counter/"):]
	metricData := strings.Split(metricValue, "/")
	metricName := metricData[0]
	value, _ := strconv.ParseInt(metricData[1], 10, 64)
	ms.AddCounter(metricName, counter(value))
}
