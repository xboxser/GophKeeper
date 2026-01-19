package model

import "io"

type FileAdd struct {
	FileName string `json:"file_path"`
	UserID   int
	Size     int64
	Body     io.ReadCloser
}

type FileAPI struct {
	Name string `json:"file_name"`
	Size int64  `json:"file_size"`
}

type File struct {
	ID     int
	UserID int
	Name   string
	Size   int64
}
