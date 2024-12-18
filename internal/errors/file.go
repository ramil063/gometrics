package errors

import (
	"fmt"
	"time"
)

type FileError struct {
	Time time.Time
	Err  error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%v %v", e.Time.Format("2006-01-02 15:04:05"), e.Err)
}

func (e *FileError) Unwrap() error {
	return e.Err
}

// NewFileError записывает ошибку err в тип FileError c текущим временем.
func NewFileError(err error) error {
	return &FileError{
		Time: time.Now(),
		Err:  err,
	}
}
