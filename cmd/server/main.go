package main

import (
	"context"
	"fmt"
	"gophkeeper/internal/config"
	"gophkeeper/internal/server/db"
	"gophkeeper/internal/server/handler"
	"gophkeeper/internal/server/handler/middleware"
	"gophkeeper/internal/server/handler/route"
	"gophkeeper/internal/server/repository"
	"gophkeeper/internal/server/service"
	"log"
	"net/http"
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

	// Получаем все объекты для User - пользователи
	userRepository := repository.NewUserDB(db)
	userService := service.NewUserService(userRepository, jwtTokenService, sessionService)
	userHandler := route.NewUserHandler(userService)

	// формируем middleware для обработки нужных запросов
	tokenMiddleware := middleware.NewTokenMiddleware(userService, jwtTokenService, sessionService)

	// Получаем все объекты для Credential - учетные данные пользователя
	credentialRepository := repository.NewCredentialDB(db)
	credentialService := service.NewCredentialService(credentialRepository)
	credentialHandler := route.NewCredentialHandler(credentialService, tokenMiddleware)

	cardRepository := repository.NewCardDB(db)
	cardService := service.NewCardService(cardRepository)
	cardHandler := route.NewCardHandler(cardService, tokenMiddleware)

	chiHandler := handler.NewChiHandler()
	// добавляем созданные роуты  в основной обработчики
	chiHandler.AddRoutes(userHandler)
	chiHandler.AddRoutes(credentialHandler)
	chiHandler.AddRoutes(cardHandler)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: chiHandler.Router,
	}

	// TODO запустить сервер в отдельном потоке
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
