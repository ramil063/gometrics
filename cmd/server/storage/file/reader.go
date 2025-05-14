package file

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ramil063/gometrics/internal/errors"
)

type Reader struct {
	file *os.File
	// добавляем reader в Reader
	reader *bufio.Reader
}

func NewReader(filename string) (*Reader, error) {
	if _, err := os.Stat(filepath.Dir(filename)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			return nil, errors.NewFileError(err)
		}
	}

	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		file, err = retryOpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666, errors.TriesTimes)
		if err != nil {
			return nil, errors.NewFileError(err)
		}
	}

	return &Reader{
		file: file,
		// создаём новый Reader
		reader: bufio.NewReader(file),
	}, nil
}

// ReadMetricsFromFile чтение метрик из файла
func ReadMetricsFromFile(filepath string) (*FStorage, error) {
	Reader, err := NewReader(filepath)
	if err != nil {
		return nil, errors.NewFileError(err)
	}
	defer Reader.Close()

	metrics, err := Reader.ReadMetrics()
	if err != nil {
		return nil, errors.NewFileError(err)
	}
	return metrics, nil
}

// ReadMetrics чтение метрик
func (r *Reader) ReadMetrics() (*FStorage, error) {
	// читаем данные до символа переноса строки
	data, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, errors.NewFileError(err)
	}

	// преобразуем данные из JSON-представления в структуру
	metrics := FStorage{}
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return nil, errors.NewFileError(err)
	}

	return &metrics, nil
}

// Close закрытие соединения с файлом
func (r *Reader) Close() error {
	// закрываем файл
	err := r.file.Close()
	if err != nil {
		return errors.NewFileError(err)
	}
	return nil
}
