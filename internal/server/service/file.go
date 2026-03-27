package service

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
	"os"
)

type FileService interface {
	AddFile(ctx context.Context, file model.FileAdd) error
	DownloadFile(ctx context.Context, file model.File) (*os.File, os.FileInfo, error)
	DeleteFile(ctx context.Context, file model.File) error

	GetFileForName(ctx context.Context, userID int, fileName string) (model.File, error)
	ListFile(ctx context.Context, userID int) ([]model.FileAPI, error)
	// StopDownloadFile - устанавливаем статус ok, т.к. скачивание прекратилось
	StopDownloadFile(ctx context.Context, file model.File) error
}

type fileService struct {
	FileRepository repository.FileRepository
}

func NewFileService(fileRepository repository.FileRepository) *fileService {
	return &fileService{
		FileRepository: fileRepository,
	}
}

func (fs *fileService) GetFileForName(ctx context.Context, userID int, fileName string) (model.File, error) {
	return fs.FileRepository.GetFileForName(ctx, userID, fileName)
}

func (fs *fileService) AddFile(ctx context.Context, file model.FileAdd) error {
	return fs.FileRepository.AddFile(ctx, file)
}

func (fs *fileService) ListFile(ctx context.Context, userID int) ([]model.FileAPI, error) {
	return fs.FileRepository.GetList(ctx, userID)
}

func (fs *fileService) DownloadFile(ctx context.Context, file model.File) (*os.File, os.FileInfo, error) {
	// Проверка файла
	f, err := os.Open(file.FilePath)
	if err != nil {
		return nil, nil, model.ErrFileNotFound
	}

	// Устанавливаем статус что файл скачивает
	err = fs.FileRepository.DownloadFile(ctx, file.UserID, file.Name)
	if err != nil {
		defer f.Close()
		return nil, nil, model.ErrFileUpdate
	}

	stat, err := f.Stat()
	if err != nil {
		defer f.Close()
		return nil, nil, err
	}

	return f, stat, nil
}

func (fs *fileService) DeleteFile(ctx context.Context, file model.File) error {
	return fs.FileRepository.DeleteFile(ctx, file.UserID, file.Name)
}

func (fs *fileService) StopDownloadFile(ctx context.Context, file model.File) error {
	return fs.FileRepository.DeleteBlockFile(ctx, file.UserID, file.Name)
}
