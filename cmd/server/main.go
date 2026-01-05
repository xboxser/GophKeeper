package main

import (
	"context"
	"fmt"
	"gophkeeper/internal/config"
	"gophkeeper/internal/server/db"
	"gophkeeper/internal/server/handler"
	"gophkeeper/internal/server/handler/route"
	"gophkeeper/internal/server/repository"
	"gophkeeper/internal/server/service"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewServerConfig()
	fmt.Println("Start server", cfg.ServerAddress)

	ctx := context.Background()
	db, err := db.NewDBPgx(ctx, *cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Получаем все объекты для Credential - учетные данные пользователя
	credentialRepository := repository.NewCredentialBD(db)
	credentialService := service.NewCredentialService(credentialRepository)
	credentialHandler := route.NewCredentialHandler(credentialService)

	paymentCardRepository := repository.NewPaymentCardBD(db)
	paymentCardService := service.NewPaymentCardService(paymentCardRepository)
	paymentCardHandler := route.NewPaymentCardHandler(paymentCardService)

	chiHandler := handler.NewChiHandler()
	chiHandler.AddRoutes(credentialHandler)
	chiHandler.AddRoutes(paymentCardHandler)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: chiHandler.Router,
	}

	fmt.Println("HTTP server started on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
