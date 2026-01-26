package service

import (
	"context"
	"errors"
	"gophkeeper/internal/model"
	mock_rep "gophkeeper/mocks/server/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_rep.NewMockCredentialRepository(ctrl)

	credentialsMock := []model.CredentialAPI{
		{
			Login:    "login",
			Password: []byte("password"),
		},
	}
	userID := 1

	mockRepo.EXPECT().GetCredentials(gomock.Any(), userID).Return(credentialsMock, nil)
	credentialService := NewCredentialService(mockRepo)

	credentials, err := credentialService.GetCredentials(context.TODO(), userID)

	require.Equal(t, credentials, credentialsMock)
	require.NoError(t, err)
}

func TestAddCredential(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name             string
		credential       model.CredentialAPI
		userID           int
		expectError      error
		credentialGetErr error
		credentialGet    model.CredentialAPI
	}{
		{
			name:        "valid 1",
			credential:  model.CredentialAPI{Login: "login", Password: []byte("password")},
			userID:      1,
			expectError: nil,
		},
		{
			name:        "valid 2",
			credential:  model.CredentialAPI{Login: "login2", Password: []byte("password2")},
			userID:      2,
			expectError: nil,
		},
		{
			name:             "error duplicate credential",
			credential:       model.CredentialAPI{Login: "login3", Password: []byte("password2")},
			userID:           123,
			expectError:      model.ErrCredentialDuplicate,
			credentialGetErr: nil,
			credentialGet:    model.CredentialAPI{Login: "login3", Password: []byte("password2")},
		},
		{
			name:        "error user",
			credential:  model.CredentialAPI{Login: "login4", Password: []byte("password2")},
			userID:      -123,
			expectError: errors.New("not userID"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mock_rep.NewMockCredentialRepository(ctrl)

			if tc.credentialGet.Login != "" {
				mockRepo.EXPECT().GetCredential(gomock.Any(), tc.credential.Login, tc.userID).Return(tc.credentialGet, nil)
			} else {
				mockRepo.EXPECT().GetCredential(gomock.Any(), tc.credential.Login, tc.userID).Return(model.CredentialAPI{}, tc.credentialGetErr)
				mockRepo.EXPECT().AddCredential(gomock.Any(), tc.credential, tc.userID).Return(tc.expectError)
			}

			credentialService := NewCredentialService(mockRepo)

			err := credentialService.AddCredential(context.Background(), tc.credential, tc.userID)
			if tc.expectError != nil {
				assert.ErrorIs(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestDeleteCredential(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name             string
		login            string
		userID           int
		expectError      error
		credentialGetErr error
		credentialGet    model.CredentialAPI
	}{
		{
			name:          "valid 1",
			login:         "login",
			userID:        1,
			expectError:   nil,
			credentialGet: model.CredentialAPI{Login: "login", Password: []byte("password2")},
		},
		{
			name:          "valid 2",
			login:         "login",
			userID:        2,
			expectError:   nil,
			credentialGet: model.CredentialAPI{Login: "login", Password: []byte("password2")},
		},
		{
			name:             "error not found credential",
			login:            "login3",
			userID:           123,
			expectError:      model.ErrCredentialNotFound,
			credentialGetErr: nil,
		},
		{
			name:          "error user",
			login:         "login4",
			userID:        -123,
			expectError:   errors.New("not userID"),
			credentialGet: model.CredentialAPI{Login: "login", Password: []byte("password2")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mock_rep.NewMockCredentialRepository(ctrl)

			if tc.credentialGet.Login != "" {
				mockRepo.EXPECT().GetCredential(gomock.Any(), gomock.Any(), tc.userID).Return(tc.credentialGet, nil)
				mockRepo.EXPECT().DeleteCredential(gomock.Any(), tc.login, tc.userID).Return(tc.expectError)
			} else {
				mockRepo.EXPECT().GetCredential(gomock.Any(), gomock.Any(), tc.userID).Return(model.CredentialAPI{}, tc.credentialGetErr)

			}

			credentialService := NewCredentialService(mockRepo)

			err := credentialService.DeleteCredential(context.Background(), tc.login, tc.userID)
			if tc.expectError != nil {
				assert.ErrorIs(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
