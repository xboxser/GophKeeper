package handler

import (
	"gophkeeper/internal/client/service"

	"github.com/spf13/cobra"
)

type Cli struct {
	ClientCmd *cobra.Command
}

func NewCli(credentialService service.CredentialService) *Cli {
	// DI: Repository → Service → Handlers → Commands

	clientCmd := &cobra.Command{
		Use:   "todo",
		Short: "Клиент для работы с gophkeeper",
		Long:  "Клиент для работы с gophkeeper",
		// TODO добавить версию
		Version: "1.0.0",
	}

	credentialHandler := NewCredentialHandler(credentialService)
	// Команда add
	addCmd := &cobra.Command{
		Use:   "addCredential",
		Short: "Добавить новую задачу",
		Example: `  todo addCredential --login "userName" --password "password"
  todo add --text "Позвонить маме"`,
		Run: credentialHandler.AddCredential,
	}
	addCmd.Flags().StringP("login", "l", "", "Логин (обязательный)")
	addCmd.MarkFlagRequired("login")
	addCmd.Flags().StringP("password", "p", "", "Пароль (обязательный)")
	addCmd.MarkFlagRequired("password")

	// Собираем дерево команд
	clientCmd.AddCommand(addCmd)

	return &Cli{ClientCmd: clientCmd}
}

func (cli *Cli) Run() error {
	return cli.ClientCmd.Execute()
}
