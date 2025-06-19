package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClearFileContent(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir := t.TempDir()

	tests := []struct {
		name          string
		setup         func() string // Функция подготовки теста (возвращает путь к файлу)
		expectedError bool
	}{
		{
			name: "file exists and is cleared successfully",
			setup: func() string {
				filePath := filepath.Join(tempDir, "existing.txt")
				err := os.WriteFile(filePath, []byte("test content"), 0644)
				if err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			expectedError: false,
		},
		{
			name: "file does not exist",
			setup: func() string {
				return filepath.Join(tempDir, "nonexistent.txt")
			},
			expectedError: false,
		},
		{
			name: "permission denied",
			setup: func() string {
				filePath := filepath.Join(tempDir, "protected.txt")
				err := os.WriteFile(filePath, []byte("test content"), 0444) // Только для чтения
				if err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setup()

			err := ClearFileContent(filePath)

			if (err != nil) != tt.expectedError {
				t.Errorf("ClearFileContent() error = %v, expectedError %v", err, tt.expectedError)
			}

			// Проверяем содержимое файла, если он существует и не ожидается ошибка
			if !tt.expectedError {
				if _, err := os.Stat(filePath); err == nil {
					content, err := os.ReadFile(filePath)
					if err != nil {
						t.Errorf("Failed to read file after clearing: %v", err)
					}
					if len(content) != 0 {
						t.Errorf("File was not cleared, content length = %d", len(content))
					}
				}
			}
		})
	}
}

func Test_retryOpenFile(t *testing.T) {
	// Создаем временный файл для теста
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Тестируем успешное открытие
	file, err := retryOpenFile(tmpfile.Name(), os.O_RDONLY, 0, []int{0})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if file == nil {
		t.Error("expected non-nil file")
	}
	file.Close()

	// Тестируем ошибку с несуществующим файлом
	_, err = retryOpenFile("nonexistent.txt", os.O_RDONLY, 0, []int{0, 1})
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
