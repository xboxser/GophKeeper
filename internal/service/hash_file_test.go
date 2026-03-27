package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFileMD5(t *testing.T) {
	// Подготовка тестовых данных
	testCases := []struct {
		name         string
		content      string
		expectedHash string
		expectError  bool
	}{
		{
			name:         "empty file",
			content:      "",
			expectedHash: "d41d8cd98f00b204e9800998ecf8427e",
			expectError:  false,
		},
		{
			name:         "file with content",
			content:      "hello world",
			expectedHash: "5eb63bbbe01eeed093cb22bb8f5acdc3",
			expectError:  false,
		},
		{
			name:         "file with multiline content",
			content:      "line1\nline2\nline3",
			expectedHash: "81facad50c8e6244de64a98cf4f56f77",
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_md5_*")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			// Запись тестового содержимого
			if _, err := tmpFile.WriteString(tc.content); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			result, err := GetFileMD5(tmpFile.Name())

			// Проверка результатов
			if tc.expectError {
				require.Error(t, err)
			}

			require.NoError(t, err)

			if result != tc.expectedHash {
				require.Equal(t, result, tc.expectedHash)
			}
		})
	}
}
