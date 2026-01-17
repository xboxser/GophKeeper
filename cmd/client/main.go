package main

import (
	"fmt"
	"gophkeeper/internal/client/handler"
	"gophkeeper/internal/client/handler/command"
	"gophkeeper/internal/client/repository"
	"gophkeeper/internal/client/service"
	"gophkeeper/internal/config"
	"net/http"
)

func main() {
	cfg := config.NewClientConfig()
	fmt.Println("Start client, for server address", cfg.ServerAddress)

	client := &http.Client{}
	senderService := service.NewSenderService(client, cfg.ServerAddress)

	encryptionService := service.NewEncryptionService()

	tokenRep := repository.NewTokenRepository("token")
	tokenService := service.NewTokenService(tokenRep)

	userService := service.NewUserService(senderService)
	userService.InitToken(tokenService)
	userCommand := command.NewUserHandler(userService)

	credentialService := service.NewCredentialService(senderService, encryptionService)
	credentialService.InitToken(tokenService)
	credentialCommand := command.NewCredentialHandler(credentialService, userService)

	cardService := service.NewCardService(senderService)
	cardService.InitToken(tokenService)
	cardCommand := command.NewCardHandler(cardService)

	cli := handler.NewCli(credentialService)
	cli.AddCommand(userCommand)
	cli.AddCommand(credentialCommand)
	cli.AddCommand(cardCommand)
	cli.Run()
}
