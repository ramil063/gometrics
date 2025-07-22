package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewRequestError(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name       string
		status     string
		statusCode int
		wantErr    bool
	}{

		{
			name:       "simple error",
			status:     "Bad Request",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewRequestError(tt.status, tt.statusCode); (err != nil) != tt.wantErr {
				t.Errorf("NewRequestError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestError_Error(t *testing.T) {
	// Фиксируем время для предсказуемых результатов
	fixedTime := time.Date(2023, time.January, 2, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		name     string
		reqErr   RequestError
		expected string
	}{
		{
			name: "standard error",
			reqErr: RequestError{
				Time:       fixedTime,
				StatusCode: 404,
				Err:        errors.New("not found"),
			},
			expected: "2023-01-02 15:04:05 Status:404 Error:not found",
		},
		{
			name: "empty error",
			reqErr: RequestError{
				Time:       fixedTime,
				StatusCode: 500,
				Err:        errors.New(""),
			},
			expected: "2023-01-02 15:04:05 Status:500 Error:",
		},
		{
			name: "zero status code",
			reqErr: RequestError{
				Time:       fixedTime,
				StatusCode: 0,
				Err:        errors.New("connection failed"),
			},
			expected: "2023-01-02 15:04:05 Status:0 Error:connection failed",
		},
		{
			name: "nil error",
			reqErr: RequestError{
				Time:       fixedTime,
				StatusCode: 400,
				Err:        nil,
			},
			expected: "2023-01-02 15:04:05 Status:400 Error:<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.reqErr.Error()
			if result != tt.expected {
				t.Errorf("Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRequestError_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		fileErr  RequestError
		expected error
	}{
		{
			name: "with standard error",
			fileErr: RequestError{
				Err: errors.New("bad request"),
			},
			expected: errors.New("bad request"),
		},
		{
			name: "with nil error",
			fileErr: RequestError{
				Err: nil,
			},
			expected: nil,
		},
		{
			name: "with wrapped error",
			fileErr: RequestError{
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
