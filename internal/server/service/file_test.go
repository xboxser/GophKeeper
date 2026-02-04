package service

import (
	"context"
	"errors"
	"gophkeeper/internal/model"
	"os"
	"testing"

	mock_rep "gophkeeper/mocks/server/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetFileForName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock_rep.NewMockFileRepository(ctrl)
	fileService := NewFileService(mockFileRepo)

	userID := 123
	fileName := "test.txt"
	expectedFile := model.File{
		ID:       1,
		UserID:   userID,
		Name:     fileName,
		FilePath: "/path/to/test.txt",
		Size:     1024,
		Status:   "active",
	}

	t.Run("successful get file for name", func(t *testing.T) {
		mockFileRepo.EXPECT().GetFileForName(gomock.Any(), userID, fileName).Return(expectedFile, nil)

		result, err := fileService.GetFileForName(context.Background(), userID, fileName)

		require.NoError(t, err)
		require.Equal(t, expectedFile, result)
	})

	t.Run("error getting file for name", func(t *testing.T) {
		expectedError := errors.New("database error")

		mockFileRepo.EXPECT().GetFileForName(gomock.Any(), userID, fileName).Return(model.File{}, expectedError)

		result, err := fileService.GetFileForName(context.Background(), userID, fileName)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Equal(t, model.File{}, result)
	})
}

func TestAddFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock_rep.NewMockFileRepository(ctrl)
	fileService := NewFileService(mockFileRepo)

	fileToAdd := model.FileAdd{
		UserID:   123,
		FileName: "test.txt",
		Size:     1024,
	}

	t.Run("successful add file", func(t *testing.T) {
		mockFileRepo.EXPECT().AddFile(gomock.Any(), fileToAdd).Return(nil)

		err := fileService.AddFile(context.Background(), fileToAdd)

		require.NoError(t, err)
	})

	t.Run("error adding file", func(t *testing.T) {
		expectedError := errors.New("database error")

		mockFileRepo.EXPECT().AddFile(gomock.Any(), fileToAdd).Return(expectedError)

		err := fileService.AddFile(context.Background(), fileToAdd)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
	})
}

func TestListFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock_rep.NewMockFileRepository(ctrl)
	fileService := NewFileService(mockFileRepo)

	userID := 123
	expectedFiles := []model.FileAPI{
		{

			Name:   "test1.txt",
			Size:   1024,
			Status: "active",
		},
		{
			Name:   "test2.txt",
			Size:   2048,
			Status: "active",
		},
	}

	t.Run("successful list files", func(t *testing.T) {
		mockFileRepo.EXPECT().GetList(gomock.Any(), userID).Return(expectedFiles, nil)

		result, err := fileService.ListFile(context.Background(), userID)

		require.NoError(t, err)
		require.Equal(t, expectedFiles, result)
	})

	t.Run("error listing files", func(t *testing.T) {
		expectedError := errors.New("database error")

		mockFileRepo.EXPECT().GetList(gomock.Any(), userID).Return(nil, expectedError)

		result, err := fileService.ListFile(context.Background(), userID)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
	})
}

func TestStopDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock_rep.NewMockFileRepository(ctrl)
	fileService := NewFileService(mockFileRepo)

	userID := 123
	fileName := "test.txt"
	testFile := model.File{
		ID:       1,
		UserID:   userID,
		Name:     fileName,
		FilePath: "/path/to/test.txt",
		Size:     1024,
		Status:   "active",
	}

	t.Run("successful stop download file", func(t *testing.T) {
		mockFileRepo.EXPECT().DeleteBlockFile(gomock.Any(), userID, fileName).Return(nil)

		err := fileService.StopDownloadFile(context.Background(), testFile)

		require.NoError(t, err)
	})

}

func TestDeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock_rep.NewMockFileRepository(ctrl)
	fileService := NewFileService(mockFileRepo)

	userID := 123
	fileName := "test.txt"
	testFile := model.File{
		ID:       1,
		UserID:   userID,
		Name:     fileName,
		FilePath: "/path/to/test.txt",
		Size:     1024,
		Status:   "active",
	}

	t.Run("successful delete file", func(t *testing.T) {
		mockFileRepo.EXPECT().DeleteFile(gomock.Any(), userID, fileName).Return(nil)

		err := fileService.DeleteFile(context.Background(), testFile)

		require.NoError(t, err)
	})

}

func TestDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock_rep.NewMockFileRepository(ctrl)
	fileService := NewFileService(mockFileRepo)

	userID := 123
	fileName := "test.txt"
	testFilePath := "/tmp/test_download_file.txt"
	testFile := model.File{
		ID:       1,
		UserID:   userID,
		Name:     fileName,
		FilePath: testFilePath,
		Size:     1024,
		Status:   "active",
	}

	// Создаем временный файл для теста
	tmpFile, err := os.Create(testFilePath)
	require.NoError(t, err)
	tmpFile.WriteString("test content")
	tmpFile.Close()
	defer os.Remove(testFilePath) // Удаляем после завершения теста

	t.Run("successful download file", func(t *testing.T) {
		mockFileRepo.EXPECT().DownloadFile(gomock.Any(), userID, fileName).Return(nil)

		resultFile, fileInfo, err := fileService.DownloadFile(context.Background(), testFile)

		require.NoError(t, err)
		require.NotNil(t, resultFile)
		require.NotNil(t, fileInfo)
		resultFile.Close()
	})

	t.Run("error during download file", func(t *testing.T) {
		expectedError := errors.New("download error")

		mockFileRepo.EXPECT().DownloadFile(gomock.Any(), userID, fileName).Return(expectedError)

		resultFile, fileInfo, err := fileService.DownloadFile(context.Background(), testFile)

		require.Error(t, err)
		require.Equal(t, model.ErrFileUpdate, err)
		require.Nil(t, resultFile)
		require.Nil(t, fileInfo)
	})

	t.Run("file does not exist on disk", func(t *testing.T) {
		nonexistentFile := model.File{
			ID:       1,
			UserID:   userID,
			Name:     "nonexistent.txt",
			FilePath: "/nonexistent/path/file.txt",
			Size:     1024,
			Status:   "active",
		}

		resultFile, fileInfo, err := fileService.DownloadFile(context.Background(), nonexistentFile)

		require.Error(t, err)
		require.Equal(t, model.ErrFileNotFound, err)
		require.Nil(t, resultFile)
		require.Nil(t, fileInfo)
	})
}
