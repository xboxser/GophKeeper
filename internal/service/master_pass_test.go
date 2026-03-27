package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestConvertMasterPassToHash(t *testing.T) {
	tests := []struct {
		name        string
		masterPass  string
		expectError bool
	}{
		{
			name:        "valid_password_simple",
			masterPass:  "simplepass123",
			expectError: false,
		},
		{
			name:        "valid_password_complex",
			masterPass:  "ComplexP@ssw0rd!2023",
			expectError: false,
		},
		{
			name:        "empty_password",
			masterPass:  "",
			expectError: false,
		},
		{
			name:        "long_password",
			masterPass:  "very_long_password_that_exceeds_normal_length_requirements_123456789",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := ConvertMasterPassToHash(tt.masterPass)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, hash)
				require.NotEmpty(t, hash)

				// Проверяем, что хэш действительно можно использовать для проверки
				err = bcrypt.CompareHashAndPassword(hash, []byte(tt.masterPass))
				require.NoError(t, err, "Generated hash should validate against original password")
			}
		})
	}
}

func TestValidateHash(t *testing.T) {
	password := "testPassword123"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	require.NotNil(t, hash)

	tests := []struct {
		name          string
		hash          []byte
		masterPass    string
		expectedError bool
	}{
		{
			name:          "correct_password_and_hash",
			hash:          hash,
			masterPass:    password,
			expectedError: false,
		},
		{
			name:          "incorrect_password",
			hash:          hash,
			masterPass:    "wrongPassword456",
			expectedError: true,
		},
		{
			name:          "empty_password_with_valid_hash",
			hash:          hash,
			masterPass:    "",
			expectedError: true,
		},
		{
			name:          "empty_hash_with_correct_password",
			hash:          []byte(""),
			masterPass:    password,
			expectedError: true,
		},
		{
			name:          "invalid_hash_format",
			hash:          []byte("invalid_hash_format"),
			masterPass:    password,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHash(tt.hash, tt.masterPass)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
