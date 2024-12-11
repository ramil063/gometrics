package errors

import (
	"fmt"
	"time"
)

type DbError struct {
	Time time.Time
	Err  error
}

func (e *DbError) Error() string {
	return fmt.Sprintf("%v %v", e.Time.Format("2006-01-02 15:04:05"), e.Err)
}

func (e *DbError) Unwrap() error {
	return e.Err
}

// NewDbError записывает ошибку err в тип DbError c текущим временем.
func NewDbError(err error) error {
	return &DbError{
		Time: time.Now(),
		Err:  err,
	}
}
