package service

import (
	"context"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/repository"
)

type FileService interface {
	GetFile(ctx context.Context, hash []byte) ([]byte, error)
	AddFile(ctx context.Context, file model.FileAdd) error
	ListFile(ctx context.Context, userID int) ([]model.FileAPI, error)
}

type fileService struct {
	FileRepository repository.FileRepository
}

func NewFileService(fileRepository repository.FileRepository) *fileService {
	return &fileService{
		FileRepository: fileRepository,
	}
}

func (fs *fileService) GetFile(ctx context.Context, hash []byte) ([]byte, error) {
	return fs.FileRepository.GetFile(ctx, hash)
}

func (fs *fileService) AddFile(ctx context.Context, file model.FileAdd) error {
	return fs.FileRepository.AddFile(ctx, file)
}

func (fs *fileService) ListFile(ctx context.Context, userID int) ([]model.FileAPI, error) {
	return fs.FileRepository.GetList(ctx, userID)
}
