package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ramil063/gometrics/cmd/server/handlers"
	"github.com/ramil063/gometrics/cmd/server/handlers/writers"
	"github.com/ramil063/gometrics/internal/hash"
	"github.com/ramil063/gometrics/internal/logger"
)

// CheckMethodMw middleware для проверки метода запроса
func CheckMethodMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем только POST, GET запросы
		if r.Method != http.MethodPost && r.Method != http.MethodGet {
			logger.WriteDebugLog("Incorrect method", "")
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
			logger.WriteDebugLog("got request with bad method", "method:"+r.Method)
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
			logger.WriteDebugLog("Error metric name is empty", "")
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
			logger.WriteDebugLog("Error in metric type (allowed 'gauge' or 'counter')", "got:"+r.PathValue("type"))
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
			logger.WriteDebugLog("Error in metric value", "")
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
			logger.WriteDebugLog("Error in metric data(update)", "")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if !issetCorrectValue {
			logger.WriteDebugLog("Error in metric value(update)", "")
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
			logger.WriteDebugLog("Error in metric type", "")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.PathValue("metric") == "" {
			logger.WriteDebugLog("Error in metric name", "")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GZIPMiddleware нужен для сжатия входящих и выходных данных
func GZIPMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		contentType := r.Header.Get("Content-Type")
		accept := r.Header.Get("Accept")
		applicationJSON := strings.Contains(contentType, "application/json")
		textHTML := strings.Contains(contentType, "text/html")
		acceptTextHTML := strings.Contains(accept, "text/html")

		if supportsGzip && (applicationJSON || textHTML || acceptTextHTML) {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := handlers.NewCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip && (applicationJSON || textHTML || acceptTextHTML) {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := handlers.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(ow, r)
	})
}

// CheckHashMiddleware проверка полученного и высчитанного хеша(подписи)
func CheckHashMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if handlers.HashKey != "" {
			body, _ := io.ReadAll(r.Body)

			headerHashSHA256 := r.Header.Get("HashSHA256")
			bodyHashSHA256 := hash.CreateSha256(body, handlers.HashKey)

			if headerHashSHA256 != bodyHashSHA256 {
				logger.WriteErrorLog("hash isn't correct", "HashSHA256")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			//возвращаем прочитанное тело обратно
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой добавления заголовка хеша при ответе
			hw := writers.NewHashWriter(w, body, handlers.HashKey)
			// меняем оригинальный http.ResponseWriter на новый
			w = hw
		}
		// передаём управление хендлеру
		next.ServeHTTP(w, r)
	})
}
