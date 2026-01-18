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

	// проверка успешного скачивания файла
	// при любой ошибке удаляем файл
	var hasError bool

	// Создаём файл
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	// Отложенная функция для закрытия файла и удаления при ошибке
	defer func() {
		out.Close()
		if hasError {
			os.Remove(path)
		}
	}()

	_, err = io.Copy(out, file.Body)
	if err != nil {
		hasError = true
		return err
	}

	err = f.saveAddFile(ctx, file)
	if err != nil {
		hasError = true
		return err
	}

	return nil
}

// saveAddFile - Сохраняем информацию по файлу в БД
func (f *fileRepository) saveAddFile(ctx context.Context, file model.FileAdd) error {
	query := `INSERT INTO files (user_id, name) VALUES ($1, $2)`
	_, err := f.DB.Exec(ctx, query, file.UserID, file.FileName)
	if err != nil {
		return err
	}
	return nil
}
