package repository

import (
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

func TestNewFileRepository(t *testing.T) {
	localStorage := "/tmp/"
	fileRepository := NewFileRepository(localStorage)
	require.Equal(t, localStorage, fileRepository.LocalStorage)
}
