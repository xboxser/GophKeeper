package command

import (
	"fmt"
	"gophkeeper/internal/client/service"

	"github.com/spf13/cobra"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (u *UserHandler) GetCommands() []*cobra.Command {

	RegisterCmd := &cobra.Command{
		Use:   "registration",
		Short: "Добавить новые учетные записи",
		Example: `  todo registration --login "userName" --password "password" --masterPass "password"
			Значение флага --masterPass: пароль для шифрования, его необходимо запомнить. Не подлежит восстановлению!!!`,
		Run: u.RegistrationRun,
	}
	RegisterCmd.Flags().StringP("login", "l", "", "Логин (Обязательный)")
	RegisterCmd.MarkFlagRequired("login")
	RegisterCmd.Flags().StringP("password", "p", "", "Пароль (Обязательный)")
	RegisterCmd.MarkFlagRequired("password")
	RegisterCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный. Не подлежит восстановлению!!!)")
	RegisterCmd.MarkFlagRequired("masterPass")

	LoginCmd := &cobra.Command{
		Use:     "login",
		Short:   "Авторизоваться",
		Example: `  todo login --login "userName" --password "password"`,
		Run:     u.LoginRun,
	}
	LoginCmd.Flags().StringP("login", "l", "", "Логин (обязательный)")
	LoginCmd.MarkFlagRequired("login")
	LoginCmd.Flags().StringP("password", "p", "", "Пароль (обязательный)")
	LoginCmd.MarkFlagRequired("password")

	// TODO добавить выход из учетной системы

	return []*cobra.Command{
		RegisterCmd,
		LoginCmd,
	}
}

// RegistrationRun - обработчик регистрации пользователя
func (u *UserHandler) RegistrationRun(cmd *cobra.Command, args []string) {
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
	if err != nil || password == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}

	token, err := u.UserService.Register(login, password, masterPass)

	if err != nil {
		fmt.Println("❌ Ошибка регистрации:", err)
		return
	}

	if token == "" {
		fmt.Println("❌ Ошибка регистрации:", "Не удалось получить токен")
		return
	}
	fmt.Println("✅ Регистрация прошла успешно")
	fmt.Println("Ваш токен:", token)
}

func (u *UserHandler) LoginRun(cmd *cobra.Command, args []string) {
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

	token, err := u.UserService.Login(login, password)

	if err != nil {
		fmt.Println("❌ Ошибка авторизации:", err)
		return
	}

	if token == "" {
		fmt.Println("❌ Ошибка авторизации:", "Не удалось получить токен")
		return
	}
	fmt.Println("✅ Авторизация прошла успешно")
	fmt.Println("Ваш токен:", token)

}
