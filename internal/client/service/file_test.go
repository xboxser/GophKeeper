package service

import (
	"encoding/json"
	"errors"
	"gophkeeper/internal/model"
	"gophkeeper/mocks/client/repository"
	"gophkeeper/mocks/client/service"
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

	// Создаем сервис
	svc := NewFileService(mockRepo, mockSender, mockEncryption)
	svc.InitToken(mockToken)

	// Создаем временный файл для теста
	tempFile, err := os.CreateTemp("", "test_file")
	require.NoError(t, err)
	testFilePath := tempFile.Name()
	defer os.Remove(testFilePath)
	tempFile.WriteString("test content")
	tempFile.Close()

	// Получаем базовое имя файла
	fileName := filepath.Base(testFilePath)

	// Устанавливаем ожидания
	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")
	mockEncryption.EXPECT().SetMasterPass("master-pass")
	mockSender.EXPECT().SendFile(gomock.Any(), "/api/files/add", gomock.Any(), fileName).Return(
		[]byte("ok"),
		&http.Response{StatusCode: http.StatusOK},
		nil,
	)

	// Вызываем метод
	err = svc.AddFile(testFilePath, "master-pass")

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
