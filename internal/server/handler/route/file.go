package route

import (
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
	//TODO добавить update & delete
	return r
}

// Pattern - поддерживает интерфейс RouteChi
func (f *FileHandler) Pattern() string {
	return "/api/files"
}

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

	err := f.FileService.AddFile(r.Context(), file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// func saveTokenToFile(filePath string, token string) error {
// 	// Получаем директорию из полного пути к файлу
// 	dir := filepath.Dir(filePath)

// 	fmt.Println(dir)
// 	// Создаём директорию (и все родительские, если нужно)
// 	if err := os.MkdirAll(dir, 0755); err != nil {
// 		return err
// 	}

// 	// Записываем файл
// 	return os.WriteFile(filePath+"ffff", []byte(token), 0600)
// }
