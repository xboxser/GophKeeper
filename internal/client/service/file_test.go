package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"gophkeeper/internal/model"
	"gophkeeper/mocks/client/repository"
	"gophkeeper/mocks/client/service"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAddFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockRepo := repository.NewMockFileRepository(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	svc := NewFileService(mockRepo, mockSender, mockEncryption)
	svc.InitToken(mockToken)

	tempFile, err := os.CreateTemp("", "test_file")
	require.NoError(t, err)
	testFilePath := tempFile.Name()
	defer os.Remove(testFilePath)
	tempFile.WriteString("test content")
	tempFile.Close()

	fileName := filepath.Base(testFilePath)

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")
	mockEncryption.EXPECT().SetMasterPass("master-pass")
	mockSender.EXPECT().SendFile(gomock.Any(), "/api/files/add", gomock.Any(), fileName).Return(
		[]byte("ok"),
		&http.Response{StatusCode: http.StatusOK},
		nil,
	)

	err = svc.AddFile(testFilePath, "master-pass")

	require.NoError(t, err)
}

func TestDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockRepo := repository.NewMockFileRepository(ctrl)

	service := NewFileService(mockRepo, mockSender, nil)
	service.InitToken(mockToken)

	fileName := "test-file.txt"

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")

	httpResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("file content"))),
	}

	mockSender.EXPECT().SendGetFile(gomock.Any(), "/api/files/download/"+fileName).Return(httpResponse, nil)
	mockRepo.EXPECT().DownloadFile(gomock.Any(), httpResponse).Return(nil)

	err := service.DownloadFile(fileName)

	require.NoError(t, err)
}

func TestDeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockRepo := repository.NewMockFileRepository(ctrl)

	// Создаем сервис
	service := NewFileService(mockRepo, mockSender, nil)
	service.InitToken(mockToken)

	fileName := "test-file.txt"

	// Подготавливаем ожидания для успешного сценария
	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")
	mockSender.EXPECT().SendDelete(gomock.Any(), "/api/files/"+fileName).Return(
		[]byte("success"),
		&http.Response{StatusCode: http.StatusNoContent},
		nil,
	)

	// Вызываем метод
	err := service.DeleteFile(fileName)

	// Проверяем результат
	require.NoError(t, err)
}
func TestListFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockRepo := repository.NewMockFileRepository(ctrl)

	// Создаем сервис
	service := NewFileService(mockRepo, mockSender, nil)
	service.InitToken(mockToken)

	t.Run("successful list files", func(t *testing.T) {
		// Подготавливаем данные
		expectedFiles := []model.FileAPI{
			{Name: "file1.txt", Size: 100},
			{Name: "file2.pdf", Size: 200},
		}
		jsonResponse, _ := json.Marshal(expectedFiles)

		// Ожидаем вызовы моков
		mockToken.EXPECT().GetToken().Return("token", nil)
		mockSender.EXPECT().SetToken("token")
		mockSender.EXPECT().SendGet(gomock.Any(), "/api/files/list").Return(
			jsonResponse,
			&http.Response{StatusCode: http.StatusOK},
			nil,
		)

		// Вызываем метод
		result, err := service.ListFile()

		// Проверяем результат
		require.NoError(t, err)
		require.Equal(t, expectedFiles, result)
	})

	t.Run("error getting token", func(t *testing.T) {
		mockToken.EXPECT().GetToken().Return("", errors.New("token error"))

		_, err := service.ListFile()

		require.Error(t, err)
		require.Contains(t, err.Error(), "token error")
	})

	t.Run("error from sender", func(t *testing.T) {
		mockToken.EXPECT().GetToken().Return("token", nil)
		mockSender.EXPECT().SetToken("token")
		mockSender.EXPECT().SendGet(gomock.Any(), "/api/files/list").Return(
			nil, nil, errors.New("sender error"),
		)

		_, err := service.ListFile()

		require.Error(t, err)
		require.Contains(t, err.Error(), "sender error")
	})

	t.Run("non-success status code", func(t *testing.T) {
		mockToken.EXPECT().GetToken().Return("token", nil)
		mockSender.EXPECT().SetToken("token")
		mockSender.EXPECT().SendGet(gomock.Any(), "/api/files/list").Return(
			[]byte("server error"),
			&http.Response{StatusCode: http.StatusInternalServerError},
			nil,
		)

		_, err := service.ListFile()

		require.Error(t, err)
		require.Contains(t, err.Error(), "error status get list files")
	})
}
