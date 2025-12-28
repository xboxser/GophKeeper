package main

import (
	"fmt"
	"gophkeeper/internal/client/handler"
	"gophkeeper/internal/client/repository"
	"gophkeeper/internal/client/service"
	"gophkeeper/internal/config"
)

func main() {
	fmt.Println("Start client")
	cfg := config.NewClientConfig()
	fmt.Println("Start client, for server address", cfg.ServerAddress)

	credentialRep := repository.NewCredentialMem()
	credentialService := service.NewCredentialService(credentialRep)

	cli := handler.NewCli(credentialService)
	cli.Run()
}
