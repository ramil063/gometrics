package server

import (
	"bytes"
	"fmt"
	"net/http/httptest"
)

func ExampleHome() {
	mockStorage := NewMemStorage()
	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/", nil)
	// Создаем Recorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	Home(rr, req, mockStorage)

	fmt.Println("Content-Type:", rr.Header().Get("Content-Type"))
	fmt.Println("Status:", rr.Code)
	fmt.Println("Body Length:", len(rr.Body.String()))
	fmt.Println(rr.Body.String()[:109])

	// Output:
	// Content-Type: text/html
	// Status: 200
	// Body Length: 185
	//
	//<!DOCTYPE html>
	//<html lang="ru">
	//<head>
	//     <meta charset="UTF-8">
	//     <title>Все метрики</title>
}

func ExampleGetValue() {
	mockStorage := NewMemStorage()
	req := httptest.NewRequest("GET", "/value/counter/met1", nil)
	rr := httptest.NewRecorder()
	GetValue(rr, req, mockStorage)
	fmt.Println("Status:", rr.Code)

	// Output:
	// Status: 200
}

func ExampleUpdate() {
	mockStorage := NewMemStorage()
	req := httptest.NewRequest("PUT", "/value/counter/met1", nil)
	rr := httptest.NewRecorder()

	Update(rr, req, mockStorage)

	fmt.Println("Status:", rr.Code)
	fmt.Println("Content-Type:", rr.Header().Get("Content-Type"))

	// Output:
	// Status: 200
	// Content-Type: text/plain; charset=utf-8
}

func ExampleGetValueMetricsJSON() {
	mockStorage := NewMemStorage()
	jsonStr := `{"id": "metric1", "type": "counter"}`
	body := []byte(jsonStr)
	body2 := bytes.NewReader(body)

	req := httptest.NewRequest("POST", "/value", body2)
	rr := httptest.NewRecorder()
	GetValueMetricsJSON(rr, req, mockStorage)

	fmt.Println("Content-Type:", rr.Header().Get("Content-Type"))
	fmt.Println("Status:", rr.Code)
	fmt.Println("Body Length:", len(rr.Body.String()))
	fmt.Println(rr.Body.String())

	// Output:
	// Content-Type: application/json
	// Status: 200
	// Body Length: 44
	//{"id":"metric1","type":"counter","delta":0}
}

func ExampleUpdateMetricsJSON() {
	mockStorage := NewMemStorage()

	jsonStr := `{"id": "metric1", "type": "counter", "delta": 1}`
	body := []byte(jsonStr) // Преобразуем строку в []byte
	body2 := bytes.NewReader(body)

	req := httptest.NewRequest("POST", "/update", body2)
	rr := httptest.NewRecorder()
	UpdateMetricsJSON(rr, req, mockStorage)

	fmt.Println("Content-Type:", rr.Header().Get("Content-Type"))
	fmt.Println("Status:", rr.Code)
	fmt.Println("Body Length:", len(rr.Body.String()))
	fmt.Println(rr.Body.String())

	// Output:
	// Content-Type: application/json
	// Status: 200
	// Body Length: 44
	//{"id":"metric1","type":"counter","delta":1}
}

func ExampleUpdates() {
	mockStorage := NewMemStorage()

	jsonStr := `[{"id": "metric1", "type": "counter", "delta": 1}]`
	body := []byte(jsonStr) // Преобразуем строку в []byte
	body2 := bytes.NewReader(body)

	req := httptest.NewRequest("POST", "/updates", body2)
	rr := httptest.NewRecorder()
	Updates(rr, req, mockStorage)

	fmt.Println("Content-Type:", rr.Header().Get("Content-Type"))
	fmt.Println("Status:", rr.Code)
	fmt.Println("Body Length:", len(rr.Body.String()))
	fmt.Println(rr.Body.String())

	// Output:
	// Content-Type: application/json
	// Status: 200
	// Body Length: 46
	//[{"id":"metric1","type":"counter","delta":1}]
}
