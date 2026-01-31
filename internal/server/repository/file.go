package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gophkeeper/internal/model"
	"gophkeeper/internal/server/db"
)

type FileRepository interface {
	// AddFile - добавляет файл в хранилище
	AddFile(ctx context.Context, file model.FileAdd) error

	// DeleteAllBlocks - удаляет все блокировки файлов
	DeleteAllBlocks(ctx context.Context) error
	DeleteBlockFile(ctx context.Context, userID int, fileName string) error
	DeleteFile(ctx context.Context, userID int, fileName string) error

	// GetFileForName - возвращает информацию по файлу
	GetFileForName(ctx context.Context, userID int, fileName string) (model.File, error)

	// GetList - возвращает список файлов
	GetList(ctx context.Context, userID int) ([]model.FileAPI, error)
}

type fileRepository struct {
	DB            db.DB
	UploadPath    string
	UploadPathTmp string
}

func NewFileRepository(db db.DB, uploadPath string) *fileRepository {
	return &fileRepository{
		DB:            db,
		UploadPath:    uploadPath,
		UploadPathTmp: uploadPath + "tmp/",
	}
}

// AddFile - добавляет файл в БД
func (f *fileRepository) AddFile(ctx context.Context, file model.FileAdd) error {
	pathTmp := f.getFilePathTmp(file.UserID, file.FileName)
	path := f.getFilePath(file.UserID, file.FileName)

	err := createDir(pathTmp)
	if err != nil {
		return err
	}

	err = createDir(path)
	if err != nil {
		return err
	}

	err = f.saveBeforeAddFile(ctx, &file)
	if err != nil {
		return errors.Join(errors.New("err before save"), err)
	}

	// проверка успешного скачивания файла
	// при любой ошибке удаляем файл
	var hasError bool

	// Создаём файл
	out, err := os.Create(pathTmp)
	if err != nil {
		return err
	}

	// Отложенная функция для закрытия файла и удаления при ошибке
	defer func() {
		out.Close()
		if hasError {
			os.Remove(pathTmp)
			ctxRemove, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			f.DeleteBlockFile(ctxRemove, file.UserID, file.FileName)
		} else {
			f.moveTmpFile(pathTmp, path)
		}
	}()

	written, err := io.Copy(out, file.Body)
	if err != nil {
		hasError = true
		return err
	}

	file.Size = written

	err = f.saveAfterAddFile(ctx, &file)
	if err != nil {
		hasError = true
		return err
	}

	return nil
}

// DeleteAllBlocks - убрать все блокировки
func (f *fileRepository) DeleteAllBlocks(ctx context.Context) error {
	// Убираем блокировки при запуске сервера
	query := `UPDATE files SET status = $1, uploaded_at=CURRENT_TIMESTAMP WHERE status != $1 AND status != $2`
	_, err := f.DB.Exec(ctx, query, model.StatusFileOk, model.StatusFileAdd)
	if err != nil {
		return err
	}

	// Для файлов которые начали загружаться, но не были сохранены
	query = `DELETE FROM files WHERE status = $1`
	_, err = f.DB.Exec(ctx, query, model.StatusFileAdd)
	if err != nil {
		return err
	}

	// удаляем папку и все файлы
	info, err := os.Stat(f.UploadPathTmp)
	if os.IsNotExist(err) {
		return nil
	}
	if info.IsDir() {
		err = os.RemoveAll(f.UploadPathTmp)
	}

	return err
}

// DeleteFile - удаляет файл пользователя
func (f *fileRepository) DeleteFile(ctx context.Context, userID int, fileName string) error {
	query := `DELETE FROM files WHERE status = $1 AND user_id = $2 AND name = $3`
	_, err := f.DB.Exec(ctx, query, model.StatusFileOk, userID, fileName)
	if err != nil {
		return err
	}
	err = os.RemoveAll(f.getFilePath(userID, fileName))
	return err
}

// DeleteAllBlocks - убрать все блокировки
func (f *fileRepository) DeleteBlockFile(ctx context.Context, userID int, fileName string) error {
	query := `UPDATE files SET status = $1, uploaded_at=CURRENT_TIMESTAMP WHERE status != $1 AND status != $4 AND user_id = $2 AND name = $3`
	_, err := f.DB.Exec(ctx, query, model.StatusFileOk, userID, fileName, model.StatusFileAdd)
	if err != nil {
		return err
	}
	// Для файлов которые начали загружаться, но не были сохранены
	query = `DELETE FROM files WHERE status = $1 AND user_id = $2 AND name = $3`
	_, err = f.DB.Exec(ctx, query, model.StatusFileAdd, userID, fileName)
	if err != nil {
		return err
	}

	return err
}

// GetList - возвращает список файлов пользователя
func (f *fileRepository) GetList(ctx context.Context, userID int) ([]model.FileAPI, error) {
	files := []model.FileAPI{}

	rows, err := f.DB.Query(ctx, "SELECT name, size, status FROM files WHERE user_id = $1 ORDER BY id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		file := model.FileAPI{}
		err := rows.Scan(&file.Name, &file.Size, &file.Status)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// GetFileForName - возвращает информацию по файлу пользователя на основе его имени
func (f *fileRepository) GetFileForName(ctx context.Context, userID int, fileName string) (model.File, error) {
	query := `SELECT id, user_id, name, size, status FROM files WHERE user_id = $1 AND name = $2 LIMIT 1`
	rows, err := f.DB.Query(ctx, query, userID, fileName)
	if err != nil {
		return model.File{}, err
	}
	defer rows.Close()

	var file model.File
	if rows.Next() {
		err := rows.Scan(&file.ID, &file.UserID, &file.Name, &file.Size, &file.Status)
		if err != nil {
			return model.File{}, err
		}
		file.FilePath = f.getFilePath(userID, fileName)
	}

	return file, nil
}

// saveBeforeAddFile - подготавливаем запись в БД перед загрузкой файла
func (f *fileRepository) saveBeforeAddFile(ctx context.Context, file *model.FileAdd) error {
	existingFile, err := f.GetFileForName(ctx, file.UserID, file.FileName)
	if err != nil {
		return err
	}

	// Если файл уже существует и он находится в блокировке, то возвращаем ошибку
	if existingFile.ID != 0 && existingFile.Status != model.StatusFileOk {
		return model.ErrFileIsBlock
	}

	var query string

	// Если файл уже существует, то обновляем его запись
	if existingFile.ID != 0 {
		query = `UPDATE files SET status = $2, uploaded_at=CURRENT_TIMESTAMP WHERE id = $1`
		_, err = f.DB.Exec(ctx, query, existingFile.ID, model.StatusFileUpdate)
	} else {
		query = `INSERT INTO files (user_id, name, size, status) VALUES ($1, $2, $3, $4)`
		_, err = f.DB.Exec(ctx, query, file.UserID, file.FileName, file.Size, model.StatusFileAdd)
	}

	if err != nil {
		return err
	}
	return nil
}

// saveAfterAddFile - Сохраняем информацию по файлу в БД
func (f *fileRepository) saveAfterAddFile(ctx context.Context, file *model.FileAdd) error {
	existingFile, err := f.GetFileForName(ctx, file.UserID, file.FileName)
	if err != nil {
		return err
	}

	// Информация о загрузке файла должна была быть уже добавлена в saveBeforeAddFile
	if existingFile.ID != 0 {
		query := `UPDATE files SET size = $1, status = $3, uploaded_at=CURRENT_TIMESTAMP WHERE id = $2`
		_, err = f.DB.Exec(ctx, query, file.Size, existingFile.ID, model.StatusFileOk)
	} else {
		err = model.ErrFileSave
	}

	if err != nil {
		return err
	}
	return nil
}

// moveTmpFile - заменяем оригинальный файл на файл из временной папки
func (f *fileRepository) moveTmpFile(pathTmp, path string) error {
	// Удаляем старый файл
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return errors.Join(errors.New("error delete old file"), err)
	}

	// Переименовываем tmp_user.mp4 в user.mp4
	err = os.Rename(pathTmp, path)
	if err != nil {
		return errors.Join(errors.New("error rename tmp file"), err)
	}
	return nil
}

func (f *fileRepository) getFilePathTmp(userID int, fileName string) string {
	return fmt.Sprintf("%s%d/%s", f.UploadPathTmp, userID, fileName)
}

// getFilePath - получить путь до файла
func (f *fileRepository) getFilePath(userID int, fileName string) string {
	return fmt.Sprintf("%s%d/%s", f.UploadPath, userID, fileName)
}

// fileExists - проверяет существование файла и что он не является каталогом
func (_ *fileRepository) fileExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

func createDir(path string) error {
	// Создаем папку
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return nil
}
