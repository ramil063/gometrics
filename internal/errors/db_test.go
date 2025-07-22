package errors

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestDBError_Error(t *testing.T) {
	// Фиксируем время для предсказуемого результата
	fixedTime := time.Date(2023, time.January, 2, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		name     string
		fileErr  DBError
		expected string
	}{
		{
			name: "simple error",
			fileErr: DBError{
				Time: fixedTime,
				Err:  errors.New("db not found"),
			},
			expected: "2023-01-02 15:04:05 db not found",
		},
		{
			name: "empty error",
			fileErr: DBError{
				Time: fixedTime,
				Err:  errors.New(""),
			},
			expected: "2023-01-02 15:04:05 ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fileErr.Error()
			if result != tt.expected {
				t.Errorf("Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDBError_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		fileErr  DBError
		expected error
	}{
		{
			name: "with standard error",
			fileErr: DBError{
				Err: errors.New("db not found"),
			},
			expected: errors.New("db not found"),
		},
		{
			name: "with nil error",
			fileErr: DBError{
				Err: nil,
			},
			expected: nil,
		},
		{
			name: "with wrapped error",
			fileErr: DBError{
				Err: fmt.Errorf("wrapper: %w", errors.New("original error")),
			},
			expected: fmt.Errorf("wrapper: %w", errors.New("original error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fileErr.Unwrap()
			// Сравниваем ошибки через errors.Is для поддержки wrapped errors
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Unwrap() = %v, want nil", result)
				}
			} else if result.Error() != tt.expected.Error() {
				t.Errorf("Unwrap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewDBError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "simple error",
			err:     errors.New("db not found"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewDBError(tt.err); (err != nil) != tt.wantErr {
				t.Errorf("NewDBError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
