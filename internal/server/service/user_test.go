package service

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
	mock_rep "gophkeeper/mocks/server/repository"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestGetUserForID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_rep.NewMockUserRepository(ctrl)

	userMock := model.User{
		ID:       1123,
		Login:    "test",
		Password: "password",
		Code:     []byte("code"),
	}
	mockRepo.EXPECT().GetUserForID(gomock.Any(), userMock.ID).Return(userMock, nil)

	userService := NewUserService(mockRepo, nil, nil)

	user, err := userService.GetUserForID(context.TODO(), userMock.ID)

	require.Equal(t, user, userMock)
	require.NoError(t, err)
}

func TestGetUserForLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_rep.NewMockUserRepository(ctrl)
	apiUser := model.APIUser{
		Login:    "test",
		Password: "password",
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(apiUser.Password), bcrypt.DefaultCost)

	userMock := model.User{
		ID:       1123,
		Login:    "test",
		Password: string(hash),
		Code:     []byte("code"),
	}
	mockRepo.EXPECT().GetUserForLogin(gomock.Any(), userMock.Login).Return(userMock, nil)

	sessionRepository := repository.NewSessionMemory()
	sessionService := NewSessionService(sessionRepository, time.Duration(1*time.Hour))

	jwtTokenService := NewTokenService("123", time.Duration(1*time.Hour))

	userService := NewUserService(mockRepo, jwtTokenService, sessionService)
	jwt, err := userService.Login(context.TODO(), apiUser)

	require.NotNil(t, jwt)
	require.NoError(t, err)
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_rep.NewMockUserRepository(ctrl)
	apiUser := model.APIUser{
		Login:    "test",
		Password: "password",
		Code:     []byte("ff"),
	}
	// Ожидаем вызов GetUserForLogin
	mockRepo.EXPECT().GetUserForLogin(gomock.Any(), apiUser.Login).Return(model.User{}, nil)

	// Ожидаем вызов RegisterUser с конкретными параметрами
	mockRepo.EXPECT().RegisterUser(gomock.Any(), apiUser).Return(1, nil)

	sessionRepository := repository.NewSessionMemory()
	sessionService := NewSessionService(sessionRepository, time.Duration(1*time.Hour))

	jwtTokenService := NewTokenService("123", time.Duration(1*time.Hour))

	userService := NewUserService(mockRepo, jwtTokenService, sessionService)
	jwt, err := userService.Register(context.TODO(), apiUser)

	require.NotNil(t, jwt)
	require.NoError(t, err)
}
