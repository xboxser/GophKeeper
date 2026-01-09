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
		Run:     ch.addCredentialRun,
	}
	addCredentialCmd.Flags().StringP("login", "l", "", "Логин (обязательный)")
	addCredentialCmd.MarkFlagRequired("login")
	addCredentialCmd.Flags().StringP("password", "p", "", "Пароль (обязательный)")
	addCredentialCmd.MarkFlagRequired("password")

	getCredentialCmd := &cobra.Command{
		Use:     "getCredential",
		Short:   "Получить список учетных данных",
		Example: ``,
		Run:     ch.getCredentialRun,
	}

	return []*cobra.Command{
		addCredentialCmd,
		getCredentialCmd,
	}
}

func (ch *CredentialHandler) addCredentialRun(cmd *cobra.Command, args []string) {
	login, err := cmd.Flags().GetString("login")
	if err != nil || login == "" {
		fmt.Println("❌ Ошибка: укажите логин (--login или -l)")
		return
	}
	password, err := cmd.Flags().GetString("password")
	if err != nil || password == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--password или -p)")
		return
	}

	err = ch.CredentialService.AddCredential(login, password)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	fmt.Printf("✅ Запись успешно добавлена\n")
}

func (ch *CredentialHandler) getCredentialRun(cmd *cobra.Command, args []string) {

	credentials, err := ch.CredentialService.GetCredentials()

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	if len(credentials) == 0 {
		fmt.Println("📦 Нет сохраненных учетных данных")
		return
	}

	fmt.Printf("🔐 Найдено %d учетных записей:\n\n", len(credentials))

	for i, cred := range credentials {
		fmt.Printf("📋 Запись #%d\n", i+1)
		fmt.Printf("   Логин: %s\n", cred.Login)
		fmt.Printf("   Пароль: %s\n", cred.Password)
		fmt.Println()
	}
}
