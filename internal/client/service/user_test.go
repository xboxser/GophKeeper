package service

import (
	"encoding/json"
	"errors"
	"gophkeeper/internal/model"
	generalService "gophkeeper/internal/service"
	"gophkeeper/mocks/client/service"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)

	userService := NewUserService(mockSender)
	userService.InitToken(mockToken)

	login := "test_user"
	password := "test_password"

	expectedUser := model.APIUser{
		Login:    login,
		Password: password,
	}

	jsonData, _ := json.Marshal(expectedUser)

	mockToken.EXPECT().SetToken("auth_token")
	mockSender.EXPECT().SendPost(gomock.Any(), "/api/user/login", jsonData).Return(
		[]byte("success"),
		&http.Response{
			StatusCode: http.StatusOK,
			Header:     map[string][]string{"Authorization": {"auth_token"}},
		},
		nil,
	)

	token, err := userService.Login(login, password)

	require.NoError(t, err)
	require.Equal(t, "auth_token", token)
}

func TestUserMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)

	userService := NewUserService(mockSender)
	userService.InitToken(mockToken)

	masterPass := "correct_master_password"

	// Создаем хешированный мастер-пароль
	hashedMasterPass, _ := generalService.ConvertMasterPassToHash(masterPass)

	expectedCode := model.Code{
		Hash: hashedMasterPass,
	}

	jsonData, _ := json.Marshal(expectedCode)

	mockToken.EXPECT().GetToken().Return("auth_token", nil)
	mockSender.EXPECT().SetToken("auth_token")
	mockSender.EXPECT().SendGet(gomock.Any(), "/api/user/code").Return(
		jsonData,
		&http.Response{
			StatusCode: http.StatusOK,
		},
		nil,
	)

	err := userService.Master(masterPass)

	require.NoError(t, err)
}

func TestUserMaster_GetTokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)

	userService := NewUserService(mockSender)
	userService.InitToken(mockToken)

	masterPass := "correct_master_password"

	mockToken.EXPECT().GetToken().Return("", errors.New("token not found"))

	err := userService.Master(masterPass)

	require.Error(t, err)
	require.Contains(t, err.Error(), "token not found")
}

func TestUserRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)

	userService := NewUserService(mockSender)
	userService.InitToken(mockToken)

	login := "test_user"
	password := "test_password"
	masterPass := "test_master_password"

	mockToken.EXPECT().SetToken("auth_token")
	mockSender.EXPECT().SendPost(gomock.Any(), "/api/user/register", gomock.Any()).Return(
		[]byte("success"),
		&http.Response{
			StatusCode: http.StatusOK,
			Header:     map[string][]string{"Authorization": {"auth_token"}},
		},
		nil,
	)

	token, err := userService.Register(login, password, masterPass)

	require.NoError(t, err)
	require.Equal(t, "auth_token", token)
}
