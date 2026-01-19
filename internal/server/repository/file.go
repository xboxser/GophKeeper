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
	GetFileForName(ctx context.Context, userID int, fileName string) (model.File, error)
	GetList(ctx context.Context, userID int) ([]model.FileAPI, error)

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

	written, err := io.Copy(out, file.Body)
	if err != nil {
		hasError = true
		return err
	}

	file.Size = written

	err = f.saveAddFile(ctx, file)
	if err != nil {
		hasError = true
		return err
	}

	return nil
}

func (f *fileRepository) GetList(ctx context.Context, userID int) ([]model.FileAPI, error) {
	files := []model.FileAPI{}

	rows, err := f.DB.Query(ctx, "SELECT name, size FROM files WHERE user_id = $1 ORDER BY id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		file := model.FileAPI{}
		err := rows.Scan(&file.Name, &file.Size)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (f *fileRepository) GetFileForName(ctx context.Context, userID int, fileName string) (model.File, error) {
	query := `SELECT id, user_id, name, size FROM files WHERE user_id = $1 AND name = $2 LIMIT 1`
	rows, err := f.DB.Query(ctx, query, userID, fileName)
	if err != nil {
		return model.File{}, err
	}
	defer rows.Close()

	var file model.File
	if rows.Next() {
		err := rows.Scan(&file.ID, &file.UserID, &file.Name, &file.Size)
		if err != nil {
			return model.File{}, err
		}
	}

	return file, nil
}

// saveAddFile - Сохраняем информацию по файлу в БД
func (f *fileRepository) saveAddFile(ctx context.Context, file model.FileAdd) error {
	existingFile, err := f.GetFileForName(ctx, file.UserID, file.FileName)
	if err != nil {
		return err
	}

	var query string

	// Если файл уже существует, то обновляем его запись
	if existingFile.ID != 0 {
		query = `UPDATE files SET size = $1, uploaded_at=CURRENT_TIMESTAMP WHERE id = $2`
		_, err = f.DB.Exec(ctx, query, file.Size, existingFile.ID)
	} else {

		query = `INSERT INTO files (user_id, name, size) VALUES ($1, $2, $3)`
		_, err = f.DB.Exec(ctx, query, file.UserID, file.FileName, file.Size)
	}

	if err != nil {
		return err
	}
	return nil
}
