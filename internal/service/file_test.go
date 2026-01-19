package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{1023, "1023 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{2048, "2.00 KB"},
		{1040000, "1015.62 KB"},
		{1048576, "1.00 MB"},
		{1073700000, "1023.96 MB"},
		{1073741824, "1.00 GB"},
		{10737418240, "10.00 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			res := FormatFileSize(tt.size)
			require.Equal(t, res, tt.expected)
		})
	}
}
