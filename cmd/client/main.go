package main

import (
	"fmt"
	"gophkeeper/internal/client/config"
	"gophkeeper/internal/client/handler"
	"gophkeeper/internal/client/handler/command"
	"gophkeeper/internal/client/repository"
	"gophkeeper/internal/client/service"
)

func main() {
	cfg := config.NewClientConfig()
	fmt.Println("Start client, for server address", cfg.ServerAddress)

	senderService := service.NewSenderService(cfg.TLS_certFile, cfg.ServerAddress)

	encryptionService := service.NewEncryptionService()

	tokenRep := repository.NewTokenRepository("token")
	tokenService := service.NewTokenService(tokenRep)

	userService := service.NewUserService(senderService)
	userService.InitToken(tokenService)
	userCommand := command.NewUserHandler(userService)

	credentialService := service.NewCredentialService(senderService, encryptionService)
	credentialService.InitToken(tokenService)
	credentialCommand := command.NewCredentialHandler(credentialService, userService)

	cardService := service.NewCardService(senderService, encryptionService)
	cardService.InitToken(tokenService)
	cardCommand := command.NewCardHandler(cardService, userService)

	fileRepository := repository.NewFileRepository(cfg.LocalStorage)
	fileService := service.NewFileService(fileRepository, senderService, encryptionService)
	fileService.InitToken(tokenService)
	fileCommand := command.NewFileHandler(fileService, userService)

	cli := handler.NewCli()
	cli.AddCommand(userCommand)
	cli.AddCommand(credentialCommand)
	cli.AddCommand(cardCommand)
	cli.AddCommand(fileCommand)
	cli.Run()
}
