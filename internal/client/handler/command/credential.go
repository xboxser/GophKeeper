package command

import (
	"fmt"
	"gophkeeper/internal/client/service"

	"github.com/spf13/cobra"
)

type CredentialHandler struct {
	CredentialService service.CredentialService
}

func NewCredentialHandler(credentialService service.CredentialService) *CredentialHandler {
	return &CredentialHandler{
		CredentialService: credentialService,
	}
}

func (ch *CredentialHandler) GetCommands() []*cobra.Command {

	addCredentialCmd := &cobra.Command{
		Use:     "addCredential",
		Short:   "Добавить новые учетные записи",
		Example: `  todo addCredential --login "userName" --password "password"`,
		Run:     ch.AddCredentialRun,
	}
	addCredentialCmd.Flags().StringP("login", "l", "", "Логин (обязательный)")
	addCredentialCmd.MarkFlagRequired("login")
	addCredentialCmd.Flags().StringP("password", "p", "", "Пароль (обязательный)")
	addCredentialCmd.MarkFlagRequired("password")

	return []*cobra.Command{
		addCredentialCmd,
	}
}

func (ch *CredentialHandler) AddCredentialRun(cmd *cobra.Command, args []string) {

	credentials, _ := ch.CredentialService.GetCredentials()
	text, _ := cmd.Flags().GetString("text")

	if text == "" {
		fmt.Println("❌ Ошибка: укажите текст задачи (--text или -t)")
		return
	}
	fmt.Printf("GetCredentials, %v\n", credentials)
	fmt.Println("text", text)
}
