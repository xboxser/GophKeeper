package main

import (
	"fmt"
	"gophkeeper/internal/config"
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

	// Получаем все объекты для Credential - учетные данные пользователя
	credentialRepository := repository.NewCredentialBD()
	credentialService := service.NewCredentialService(credentialRepository)
	credentialHandler := route.NewCredentialHandler(credentialService)

	paymentCardRepository := repository.NewPaymentCardBD()
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
