package service

import (
	"context"
	"errors"
	"gophkeeper/internal/model"
	mock_rep "gophkeeper/mocks/server/repository"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		inputCode      string
		mockSetup      func(*mock_rep.MockSessionRepository)
		expectedResult model.TokenAuth
		expectedError  error
	}{
		{
			name:      "successful_get_session",
			inputCode: "valid-code",
			mockSetup: func(repo *mock_rep.MockSessionRepository) {
				repo.EXPECT().GetSession(gomock.Any(), "valid-code").Return(
					model.Session{UserID: 123}, nil)
			},
			expectedResult: model.TokenAuth{UserID: 123},
			expectedError:  nil,
		},
		{
			name:      "error_from_repository",
			inputCode: "invalid-code",
			mockSetup: func(repo *mock_rep.MockSessionRepository) {
				repo.EXPECT().GetSession(gomock.Any(), "invalid-code").Return(
					model.Session{}, errors.New("session not found"))
			},
			expectedResult: model.TokenAuth{},
			expectedError:  errors.New("session not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_rep.NewMockSessionRepository(ctrl)
			tt.mockSetup(mockRepo)

			sessionService := NewSessionService(mockRepo, time.Hour)

			result, err := sessionService.GetSession(context.Background(), tt.inputCode)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestAddSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		inputTokenAuth model.TokenAuth
		mockSetup      func(*mock_rep.MockSessionRepository)
		expireTime     time.Duration
		expectSuccess  bool
		expectedError  error
	}{
		{
			name: "successful_add_session",
			inputTokenAuth: model.TokenAuth{
				UserID: 123,
			},
			expireTime: time.Hour,
			mockSetup: func(repo *mock_rep.MockSessionRepository) {
				repo.EXPECT().AddSession(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, session model.Session) error {
						// Проверяем, что UUID не пустой и UserID совпадает
						require.NotEmpty(t, session.UUID)
						require.Equal(t, 123, session.UserID)
						return nil
					},
				)
			},
			expectSuccess: true,
		},
		{
			name: "error_from_repository",
			inputTokenAuth: model.TokenAuth{
				UserID: 456,
			},
			expireTime: time.Hour,
			mockSetup: func(repo *mock_rep.MockSessionRepository) {
				repo.EXPECT().AddSession(gomock.Any(), gomock.Any()).Return(
					errors.New("database error"))
			},
			expectSuccess: false,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_rep.NewMockSessionRepository(ctrl)
			tt.mockSetup(mockRepo)

			sessionService := NewSessionService(mockRepo, tt.expireTime)

			result, err := sessionService.AddSession(context.Background(), tt.inputTokenAuth)

			if tt.expectSuccess {
				require.NoError(t, err)
				require.NotEmpty(t, result.UUID)
				require.Equal(t, tt.inputTokenAuth.UserID, result.UserID)

			} else {
				require.Error(t, err)
				if tt.expectedError != nil {
					require.Equal(t, tt.expectedError.Error(), err.Error())
				}
			}
		})
	}
}
