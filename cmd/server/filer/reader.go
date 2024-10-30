package filer

import (
	"bufio"
	"encoding/json"
	"github.com/ramil063/gometrics/cmd/agent/storage"
	"os"
)

type Reader struct {
	file *os.File
	// добавляем reader в Reader
	reader *bufio.Reader
}

func NewReader(filename string) (*Reader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Reader{
		file: file,
		// создаём новый Reader
		reader: bufio.NewReader(file),
	}, nil
}

func (r *Reader) ReadMonitor() (*storage.Monitor, error) {
	// читаем данные до символа переноса строки
	data, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// преобразуем данные из JSON-представления в структуру
	monitor := storage.Monitor{}
	err = json.Unmarshal(data, &monitor)
	if err != nil {
		return nil, err
	}

	return &monitor, nil
}

func (r *Reader) Close() error {
	// закрываем файл
	return r.file.Close()
}
