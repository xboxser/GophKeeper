package model

import "io"

type FileAdd struct {
	FileName string `json:"file_path"`
	UserID   int
	Body     io.ReadCloser
}
