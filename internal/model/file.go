package model

import (
	"errors"
	"io"
)

var (
	ErrFileEmptyFileName = errors.New("error empty file name")
	ErrFileNotFound      = errors.New("file not found")
	// ErrFileSave - ошибка сохранения файла
	ErrFileSave = errors.New("error save file")
	// ErrFileIsBlock - файл заблокирован
	ErrFileIsBlock = errors.New("file is block")
	ErrFileUpdate  = errors.New("error file update")
)

// StatusFile - статус файла
type StatusFile string

const (
	// StatusFileAdd - Файл загружается на сервер
	StatusFileAdd StatusFile = "add"
	// StatusFileUpdate - обновляем файл
	StatusFileUpdate StatusFile = "upd."

	// StatusFileDown - файл скачивается пользователем
	StatusFileDown StatusFile = "down"

	// StatusFileOl - файл загружен и хранится на сервере
	StatusFileOk StatusFile = "ok"
)

type FileAdd struct {
	FileName string `json:"file_path"`
	UserID   int
	Size     int64
	Body     io.ReadCloser
}

// FileAPI - объект для передачи информации об файле по АПИ
type FileAPI struct {
	Name   string `json:"file_name"`
	Size   int64  `json:"file_size"`
	Status string `json:"file_status"`
}

// File - объект с информацией о файле на сервере
type File struct {
	ID       int
	UserID   int
	Name     string
	Status   StatusFile
	FilePath string
	Size     int64
}
