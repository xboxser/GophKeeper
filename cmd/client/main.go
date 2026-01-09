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

	credentialService := service.NewCredentialService(senderService)
	credentialService.InitToken(tokenService)
	credentialHandler := command.NewCredentialHandler(credentialService)

	userService := service.NewUserService(senderService)
	userService.InitToken(tokenService)
	userCommand := command.NewUserHandler(userService)

	cli := handler.NewCli(credentialService)
	cli.AddCommand(userCommand)
	cli.AddCommand(credentialHandler)
	cli.Run()
}
