// Package writers осуществляется работа с хешем
package writers

import (
	"net/http"

	hashInternal "github.com/ramil063/gometrics/internal/hash"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type hashWriter struct {
	w    http.ResponseWriter
	key  string
	body []byte
}

func NewHashWriter(w http.ResponseWriter, body []byte, key string) *hashWriter {
	return &hashWriter{
		w:    w,
		key:  key,
		body: body,
	}
}

func (hw *hashWriter) Header() http.Header {
	return hw.w.Header()
}

func (hw *hashWriter) Write(p []byte) (int, error) {
	return hw.w.Write(p)
}

func (hw *hashWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		hw.w.Header().Set("HashSHA256", hw.CreateSha256())
	}
	hw.w.WriteHeader(statusCode)
}

// CreateSha256 создаем хеш
func (hw *hashWriter) CreateSha256() string {
	return hashInternal.CreateSha256(hw.body, hw.key)
}
