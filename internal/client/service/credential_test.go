package service

import (
	"encoding/json"
	"gophkeeper/internal/model"
	"gophkeeper/mocks/client/service"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestDeleteCredential(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	service := NewCredentialService(mockSender, mockEncryption)
	service.InitToken(mockToken)

	credID := "123"

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")
	mockSender.EXPECT().SendDelete(gomock.Any(), "/api/credentials/"+credID).Return(
		[]byte("success"),
		&http.Response{StatusCode: http.StatusNoContent},
		nil,
	)

	err := service.DeleteCredential(credID)

	require.NoError(t, err)
}

func TestGetCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	service := NewCredentialService(mockSender, mockEncryption)
	service.InitToken(mockToken)

	masterPass := "master_password"

	expectedCredentialsAPI := []model.CredentialAPI{
		{
			IncID:    1,
			Login:    "test_user",
			Password: []byte("encrypted_password"),
		},
		{
			IncID:    2,
			Login:    "another_user",
			Password: []byte("another_encrypted_password"),
		},
	}

	jsonData, _ := json.Marshal(expectedCredentialsAPI)

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")
	mockSender.EXPECT().SendGet(gomock.Any(), "/api/credentials").Return(
		jsonData,
		&http.Response{StatusCode: http.StatusOK},
		nil,
	)

	mockEncryption.EXPECT().SetMasterPass(masterPass)
	mockEncryption.EXPECT().Decrypt([]byte("encrypted_password")).Return("plain_password", nil)
	mockEncryption.EXPECT().Decrypt([]byte("another_encrypted_password")).Return("another_plain_password", nil)

	result, err := service.GetCredentials(masterPass)

	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, 1, result[0].ID)
	require.Equal(t, "test_user", result[0].Login)
	require.Equal(t, "plain_password", result[0].Password)
	require.Equal(t, 2, result[1].ID)
	require.Equal(t, "another_user", result[1].Login)
	require.Equal(t, "another_plain_password", result[1].Password)
}

func TestAddCredential(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	service := NewCredentialService(mockSender, mockEncryption)
	service.InitToken(mockToken)

	login := "test_user"
	password := "test_password"
	masterPass := "master_password"

	expectedCredentialAPI := model.CredentialAPI{
		Login:    login,
		Password: []byte("encrypted_password"),
	}

	jsonData, _ := json.Marshal(expectedCredentialAPI)

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")
	mockEncryption.EXPECT().SetMasterPass(masterPass)
	mockEncryption.EXPECT().Encrypt(password).Return([]byte("encrypted_password"), nil)
	mockSender.EXPECT().SendPost(gomock.Any(), "/api/credentials", jsonData).Return(
		[]byte("success"),
		&http.Response{StatusCode: http.StatusOK},
		nil,
	)

	err := service.AddCredential(login, password, masterPass)

	require.NoError(t, err)
}

func TestUpdateCredential(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	service := NewCredentialService(mockSender, mockEncryption)
	service.InitToken(mockToken)

	login := "updated_user"
	password := "updated_password"
	masterPass := "master_password"
	id := "1"

	expectedCredentialAPI := model.CredentialAPI{
		Login:    login,
		Password: []byte("encrypted_password"),
		IncID:    1,
	}

	jsonData, _ := json.Marshal(expectedCredentialAPI)

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")
	mockEncryption.EXPECT().SetMasterPass(masterPass)
	mockEncryption.EXPECT().Encrypt(password).Return([]byte("encrypted_password"), nil)
	mockSender.EXPECT().SendPut(gomock.Any(), "/api/credentials", jsonData).Return(
		[]byte("success"),
		&http.Response{StatusCode: http.StatusOK},
		nil,
	)

	err := service.UpdateCredential(login, password, masterPass, id)

	require.NoError(t, err)
}
