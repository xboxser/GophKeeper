package repository

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFileName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid content 1",
			input:    `attachment; filename="test.txt"`,
			expected: "test.txt",
		},
		{
			name:     "No filename",
			input:    "attachment; size=12345",
			expected: "",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Invalid content",
			input:    "invalid disposition",
			expected: "",
		},
		{
			name:     "Valid content 1",
			input:    `inline; filename="document.pdf"; size=54321`,
			expected: "document.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFileName(tt.input)
			if result != tt.expected {
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewFileRepository(tempDir + "/")

	// Подготовим фейковый HTTP-ответ
	content := []byte("test file content")
	response := &http.Response{
		Body:          io.NopCloser(bytes.NewReader(content)),
		ContentLength: int64(len(content)),
		Header:        make(http.Header),
	}
	response.Header.Set("Content-Disposition", "attachment; filename=test.txt")

	err := repo.DownloadFile(context.Background(), response)

	require.NoError(t, err)

	// проверяем, что файл был создан
	filePath := filepath.Join(tempDir, "test.txt")
	require.FileExists(t, filePath)

	// проверяем содержимое файла
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	require.Equal(t, content, data)
}

func TestNewFileRepository(t *testing.T) {
	localStorage := "/tmp/"
	fileRepository := NewFileRepository(localStorage)
	require.Equal(t, localStorage, fileRepository.LocalStorage)
}
