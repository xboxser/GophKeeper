package service

import (
	"gophkeeper/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	secret := "test-secret-key"
	tokenExp := time.Hour
	tokenService := NewTokenService(secret, tokenExp)

	session := model.Session{UUID: "test-user-id"}
	validToken, err := tokenService.BuildJWTString(session)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		tokenString  string
		expectedUser string
		expectError  bool
		errorType    error
	}{
		{
			name:         "valid token",
			tokenString:  validToken,
			expectedUser: "test-user-id",
			expectError:  false,
		},
		{
			name:        "invalid token format",
			tokenString: "invalid.token.format",
			expectError: true,
			errorType:   nil,
		},
		{
			name:        "malformed token",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
			expectError: true,
			errorType:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := tokenService.GetUser(tt.tokenString)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.Equal(t, tt.errorType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, userID)
			}
		})
	}
}
