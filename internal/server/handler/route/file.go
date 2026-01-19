package route

import (
	"encoding/json"
	"fmt"
	"gophkeeper/internal/model"
	"gophkeeper/internal/server/handler/middleware"
	"gophkeeper/internal/server/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FileHandler struct {
	FileService     service.FileService
	TokenMiddleware middleware.TokenMiddleware
}

func NewFileHandler(f service.FileService, t middleware.TokenMiddleware) *FileHandler {
	return &FileHandler{
		FileService:     f,
		TokenMiddleware: t,
	}
}

// Routes - поддерживает интерфейс RouteChi
func (f *FileHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(f.TokenMiddleware.CheckToken)
	r.Post("/add", f.addFile)
	r.Get("/list", f.listFile)
	r.Get("/download/{filename}", f.downloadFile)
	//TODO добавить delete
	return r
}

// Pattern - поддерживает интерфейс RouteChi
func (f *FileHandler) Pattern() string {
	return "/api/files"
}

// addFile - обработчик добавления файла
func (f *FileHandler) addFile(w http.ResponseWriter, r *http.Request) {
	tokenAuth := f.TokenMiddleware.GetUserRequest(r)

	if tokenAuth.UserID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	filename := r.Header.Get("filename")
	if filename == "" {
		http.Error(w, "filename required", http.StatusBadRequest)
		return
	}

	file := model.FileAdd{
		UserID:   tokenAuth.UserID,
		FileName: filename,
		Body:     r.Body,
	}

	//TODO подумать над контекстом
	err := f.FileService.AddFile(r.Context(), file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// listFile - обработчик получение списка файлов
func (f *FileHandler) listFile(w http.ResponseWriter, r *http.Request) {
	tokenAuth := f.TokenMiddleware.GetUserRequest(r)

	if tokenAuth.UserID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	files, err := f.FileService.ListFile(r.Context(), tokenAuth.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(files)
}

// downloadFile - обработчик запроса на скачивание файла
func (f *FileHandler) downloadFile(w http.ResponseWriter, r *http.Request) {
	tokenAuth := f.TokenMiddleware.GetUserRequest(r)
	if tokenAuth.UserID == 0 {
		http.Error(w, "Invalid user token", http.StatusUnauthorized)
		return
	}

	filename := chi.URLParam(r, "filename")
	if filename == "" {
		http.Error(w, model.ErrFileEmptyFileName.Error(), http.StatusBadRequest)
		return
	}

	file, err := f.FileService.GetFileForName(r.Context(), tokenAuth.UserID, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if file.ID == 0 {
		http.Error(w, model.ErrFileNotFound.Error(), http.StatusNotFound)
		return
	}

	osFile, stat, err := f.FileService.DownloadFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer osFile.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", fmt.Sprint(stat.Size()))
	w.Header().Set("Cache-Control", "no-cache")

	http.ServeContent(w, r, filename, stat.ModTime(), osFile)
}
