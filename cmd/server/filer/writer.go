package filer

import (
	"bufio"
	"encoding/json"
	"github.com/ramil063/gometrics/cmd/agent/storage"
	"os"
)

type Writer struct {
	file *os.File
	// добавляем Writer в Writer
	writer *bufio.Writer
}

func NewWriter(filename string) (*Writer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Writer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (w *Writer) WriteMonitor(monitor *storage.Monitor) error {
	data, err := json.Marshal(&monitor)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := w.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := w.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return w.writer.Flush()
}

func (w *Writer) Close() error {
	return w.file.Close()
}
