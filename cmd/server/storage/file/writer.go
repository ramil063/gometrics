package file

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/ramil063/gometrics/internal/errors"
	"github.com/ramil063/gometrics/internal/logger"
)

type Writer struct {
	file *os.File
	// добавляем Writer в Writer
	writer *bufio.Writer
}

func NewWriter(filename string) (*Writer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		file, err = retryOpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666, errors.TriesTimes)
		if err != nil {
			return nil, errors.NewFileError(err)
		}
	}

	return &Writer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func WriteMetricsToFile(metrics *FStorage, filepath string) error {
	Writer, err := NewWriter(filepath)
	if err != nil {
		logger.WriteErrorLog("error create metrics writer", err.Error())
		return err
	}
	defer Writer.Close()

	err = Writer.WriteMetrics(metrics)
	if err != nil {
		logger.WriteErrorLog("error write metrics", err.Error())
		return errors.NewFileError(err)
	}
	return nil
}

func (w *Writer) WriteMetrics(metrics *FStorage) error {
	data, err := json.Marshal(&metrics)
	if err != nil {
		return errors.NewFileError(err)
	}

	// записываем событие в буфер
	if _, err := w.writer.Write(data); err != nil {
		return errors.NewFileError(err)
	}

	// добавляем перенос строки
	if err := w.writer.WriteByte('\n'); err != nil {
		return errors.NewFileError(err)
	}

	// записываем буфер в файл
	err = w.writer.Flush()
	if err != nil {
		return errors.NewFileError(err)
	}
	return nil
}

func (w *Writer) Close() error {
	err := w.file.Close()
	if err != nil {
		return errors.NewFileError(err)
	}
	return nil
}
