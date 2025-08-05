package gzip

import (
	"bytes"
	"compress/gzip"
)

// CompressData Функция для сжатия данных
func CompressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	if err = gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
