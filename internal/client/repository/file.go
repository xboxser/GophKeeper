package repository

import (
	"context"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

type FileRepository interface {
	// LoadFile - скачиваем на клиент файл с сервера
	DownloadFile(ctx context.Context, response *http.Response) error
}

type fileRepository struct {
	LocalStorage string
}

func NewFileRepository(ls string) *fileRepository {
	return &fileRepository{
		LocalStorage: ls,
	}
}

func (f *fileRepository) DownloadFile(ctx context.Context, response *http.Response) error {
	// Создаем папку
	dir := filepath.Dir(f.LocalStorage)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	fileName := getFileName(response.Header.Get("Content-Disposition"))
	if fileName == "" {
		return errors.New("file name not found")

	}
	totalSize := response.ContentLength

	out, err := os.Create(f.LocalStorage + fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	bar := pb.Full.Start64(totalSize).SetWidth(80)
	reader := bar.NewProxyReader(response.Body)

	_, err = io.Copy(out, reader)
	bar.Finish()

	return nil
}

func getFileName(disposition string) string {
	if disposition != "" {
		_, params, err := mime.ParseMediaType(disposition)
		if err == nil {
			// Извлекаем filename из header
			if remoteFilename, ok := params["filename"]; ok {
				return remoteFilename
			}
		}
	}
	return ""
}
