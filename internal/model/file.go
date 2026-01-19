package model

import (
	"errors"
	"io"
)

var (
	ErrFileEmptyFileName = errors.New("error empty file name")
	ErrFileNotFound      = errors.New("file not found")
)

type FileAdd struct {
	FileName string `json:"file_path"`
	UserID   int
	Size     int64
	Body     io.ReadCloser
}

// FileAPI - объект для передачи информации об файле по АПИ
type FileAPI struct {
	Name string `json:"file_name"`
	Size int64  `json:"file_size"`
}

// File - объект с информацией о файле на сервере
type File struct {
	ID       int
	UserID   int
	Name     string
	FilePath string
	Size     int64
}
