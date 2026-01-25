package command

import (
	"fmt"
	"gophkeeper/internal/client/service"

	"github.com/spf13/cobra"
)

type CredentialHandler struct {
	CredentialService service.CredentialService
	UserMasterService service.UserMasterService
}

func NewCredentialHandler(credentialService service.CredentialService, u service.UserMasterService) *CredentialHandler {
	return &CredentialHandler{
		CredentialService: credentialService,
		UserMasterService: u,
	}
}

func (ch *CredentialHandler) GetCommands() []*cobra.Command {

	addCredentialCmd := &cobra.Command{
		Use:     "credential-add",
		Short:   "Добавить новые учетные записи",
		Example: `  todo credential-add --login "userName" --password "password" --masterPass "password"`,
		Run:     ch.addCredentialRun,
	}
	addCredentialCmd.Flags().StringP("login", "l", "", "Логин (обязательный)")
	addCredentialCmd.MarkFlagRequired("login")
	addCredentialCmd.Flags().StringP("password", "p", "", "Пароль (обязательный)")
	addCredentialCmd.MarkFlagRequired("password")
	addCredentialCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный)")
	addCredentialCmd.MarkFlagRequired("masterPass")

	getCredentialCmd := &cobra.Command{
		Use:     "credential-get",
		Short:   "Получить список учетных данных",
		Example: `  todo credential-get --masterPass "password"`,
		Run:     ch.getCredentialRun,
	}
	getCredentialCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный)")
	getCredentialCmd.MarkFlagRequired("masterPass")

	deleteCredentialCmd := &cobra.Command{
		Use:     "credential-del",
		Short:   "Получить список учетных данных",
		Example: `  todo credential-del --masterPass "password" --login "userName"`,
		Run:     ch.deleteCredentialRun,
	}
	deleteCredentialCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный)")
	deleteCredentialCmd.MarkFlagRequired("masterPass")
	deleteCredentialCmd.Flags().StringP("login", "l", "", "Логин (обязательный)")
	deleteCredentialCmd.MarkFlagRequired("login")

	return []*cobra.Command{
		addCredentialCmd,
		getCredentialCmd,
		deleteCredentialCmd,
	}
}

func (ch *CredentialHandler) addCredentialRun(cmd *cobra.Command, _ []string) {
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
	masterPass, err := cmd.Flags().GetString("masterPass")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}

	err = ch.UserMasterService.Master(masterPass)
	if err != nil {
		fmt.Println("❌ Ошибка проверки мастер пароля:", err)
		return
	}

	err = ch.CredentialService.AddCredential(login, password, masterPass)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	fmt.Printf("✅ Запись успешно добавлена\n")
}

func (ch *CredentialHandler) deleteCredentialRun(cmd *cobra.Command, _ []string) {
	masterPass, err := cmd.Flags().GetString("masterPass")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}

	err = ch.UserMasterService.Master(masterPass)
	if err != nil {
		fmt.Println("❌ Ошибка проверки мастер пароля:", err)
		return
	}

	login, err := cmd.Flags().GetString("login")
	if err != nil || login == "" {
		fmt.Println("❌ Ошибка: укажите логин (--login или -l)")
		return
	}

	err = ch.CredentialService.DeleteCredential(login)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	fmt.Printf("✅ Запись успешно удалена\n")
}

func (ch *CredentialHandler) getCredentialRun(cmd *cobra.Command, _ []string) {
	masterPass, err := cmd.Flags().GetString("masterPass")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}

	credentials, err := ch.CredentialService.GetCredentials(masterPass)

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
		fmt.Printf("   Пароль: %s\n", string(cred.Password))
		fmt.Println()
	}
}
