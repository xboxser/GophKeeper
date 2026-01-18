package repository

import (
	"context"
	"fmt"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
	"io"
	"os"
	"path/filepath"
)

type FileRepository interface {
	GetFile(ctx context.Context, hash []byte) ([]byte, error)
	AddFile(ctx context.Context, file model.FileAdd) error
}

type fileRepository struct {
	DB         db.DB
	UploadPath string
}

func NewFileRepository(db db.DB, uploadPath string) *fileRepository {
	return &fileRepository{
		DB:         db,
		UploadPath: uploadPath,
	}
}

func (f *fileRepository) GetFile(ctx context.Context, hash []byte) ([]byte, error) {
	return []byte{}, nil
}

func (f *fileRepository) AddFile(ctx context.Context, file model.FileAdd) error {
	path := fmt.Sprintf("%s%d/%s", f.UploadPath, file.UserID, file.FileName)

	// Создаем папку
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Создаём файл
	out, err := os.Create(path)
	if err != nil {

		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file.Body)
	if err != nil {

		return err
	}
	return nil
}
