package errors

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestFileError_Error(t *testing.T) {
	// Фиксируем время для предсказуемого результата
	fixedTime := time.Date(2023, time.January, 2, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		name     string
		fileErr  FileError
		expected string
	}{
		{
			name: "simple error",
			fileErr: FileError{
				Time: fixedTime,
				Err:  errors.New("file not found"),
			},
			expected: "2023-01-02 15:04:05 file not found",
		},
		{
			name: "empty error",
			fileErr: FileError{
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

func TestFileError_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		fileErr  FileError
		expected error
	}{
		{
			name: "with standard error",
			fileErr: FileError{
				Err: errors.New("file not found"),
			},
			expected: errors.New("file not found"),
		},
		{
			name: "with nil error",
			fileErr: FileError{
				Err: nil,
			},
			expected: nil,
		},
		{
			name: "with wrapped error",
			fileErr: FileError{
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

func TestNewFileError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "simple error",
			err:     errors.New("file not found"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewFileError(tt.err); (err != nil) != tt.wantErr {
				t.Errorf("NewFileError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
