package main

import (
	"context"
	"fmt"
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/db"
	"gophkeeper/internal/server/handler"
	"gophkeeper/internal/server/handler/middleware"
	"gophkeeper/internal/server/handler/route"
	"gophkeeper/internal/server/repository"
	"gophkeeper/internal/server/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// TODO добавить логер

	cfg := config.NewServerConfig()
	fmt.Println("Start server", cfg.ServerAddress)

	ctx := context.Background()
	db, err := db.NewDBPgx(ctx, *cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	jwtTokenService := service.NewTokenService(cfg.JWT_SECRET, cfg.TokenExpiresAt)

	sessionRepository := repository.NewSessionMemory()
	sessionService := service.NewSessionService(sessionRepository, cfg.TokenExpiresAt)

	counterRepository := repository.NewUserCounterDB(db)

	// Получаем все объекты для User - пользователи
	userRepository := repository.NewUserDB(db, counterRepository)
	userService := service.NewUserService(userRepository, jwtTokenService, sessionService)

	// формируем middleware для обработки нужных запросов
	tokenMiddleware := middleware.NewTokenMiddleware(userService, jwtTokenService, sessionService)
	userHandler := route.NewUserHandler(userService, tokenMiddleware)

	// Получаем все объекты для Credential - учетные данные пользователя
	credentialRepository := repository.NewCredentialDB(db, counterRepository)
	credentialService := service.NewCredentialService(credentialRepository)
	credentialHandler := route.NewCredentialHandler(credentialService, tokenMiddleware)

	cardRepository := repository.NewCardDB(db, counterRepository)
	cardService := service.NewCardService(cardRepository)
	cardHandler := route.NewCardHandler(cardService, tokenMiddleware)

	fileRepository := repository.NewFileRepository(db, cfg.UploadPath)
	fileService := service.NewFileService(fileRepository)
	fileHandler := route.NewFileHandler(fileService, tokenMiddleware)

	chiHandler := handler.NewChiHandler()
	// добавляем созданные роуты  в основной обработчики
	chiHandler.AddRoutes(userHandler)
	chiHandler.AddRoutes(credentialHandler)
	chiHandler.AddRoutes(cardHandler)
	chiHandler.AddRoutes(fileHandler)

	chiHandler.Router.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Println("f")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Данные зашифрованы TLS", "host": "localhost"}`))
	})

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: chiHandler.Router,

		ReadHeaderTimeout: 15 * time.Second, // таймаут на чтение заголовков запроса
		ReadTimeout:       0,                // таймаут на чтение всего запроса (заголовки + тело)
		WriteTimeout:      60 * time.Minute, // таймаут на отправку ответа клиенту
		IdleTimeout:       5 * time.Minute,  // таймаут для неактивных соединений
	}

	// канал для получения сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := server.ListenAndServeTLS(cfg.TLS_certFile, cfg.TLS_keyFile); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Ждём сигнал остановки
	<-stop

	ctxStopServer, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxStopServer); err != nil {
		fmt.Printf("Ошибка при завершении сервера: %v\n", err)
		return
	}
}
