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
		Use:     "registration",
		Short:   "Добавить новые учетные записи",
		Example: `  todo registration --login "userName" --password "password"`,
		Run:     u.RegistrationRun,
	}
	RegisterCmd.Flags().StringP("login", "l", "", "Логин (обязательный)")
	RegisterCmd.MarkFlagRequired("login")
	RegisterCmd.Flags().StringP("password", "p", "", "Пароль (обязательный)")
	RegisterCmd.MarkFlagRequired("password")

	return []*cobra.Command{
		RegisterCmd,
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
		fmt.Println("❌ Ошибка: укажите логин (--password или -p)")
		return
	}

	token, err := u.UserService.Register(login, password)

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
