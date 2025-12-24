package main

import (
	"fmt"
	"gophkeeper/internal/config"
)

func main() {
	fmt.Println("Start client")
	cfg := config.NewClientConfig()
	fmt.Println("server", cfg.ServerAddress)

}
