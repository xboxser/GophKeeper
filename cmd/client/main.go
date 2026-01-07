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

	tokenRep := repository.NewTokenRepository()
	tokenService := service.NewTokenService(tokenRep)

	credentialRep := repository.NewCredentialMem()
	credentialService := service.NewCredentialService(credentialRep)
	credentialHandler := command.NewCredentialHandler(credentialService)

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, senderService)
	userService.InitToken(tokenService)
	userCommand := command.NewUserHandler(userService)

	cli := handler.NewCli(credentialService)
	cli.AddCommand(userCommand)
	cli.AddCommand(credentialHandler)
	cli.Run()
}
