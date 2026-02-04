package command

import (
	"errors"
	"fmt"
	service "gophkeeper/mocks/client/services"
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetCommands(t *testing.T) {
	mockUserService := &service.MockUserService{}
	userHandler := NewUserHandler(mockUserService)

	commands := userHandler.GetCommands()

	assert.Len(t, commands, 3)

	commandUses := make([]string, 0, 3)
	for _, cmd := range commands {
		commandUses = append(commandUses, cmd.Use)
	}
	assert.Contains(t, commandUses, "registration")
	assert.Contains(t, commandUses, "login")
	assert.Contains(t, commandUses, "master")
}

func TestMasterRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := service.NewMockUserService(ctrl)
	userHandler := NewUserHandler(mockUserService)

	tests := []struct {
		name            string
		inputMasterPass string
		masterPassValid bool
		expectedOutput  string
	}{
		{
			name:            "valid master password",
			inputMasterPass: "correctPassword",
			masterPassValid: true,
			expectedOutput:  "✅ Проверка мастер пароля прошла успешно",
		},
		{
			name:            "invalid master password",
			inputMasterPass: "wrongPassword",
			masterPassValid: false,
			expectedOutput:  "❌ Ошибка не верный мастер пароль:",
		},
		{
			name:            "empty master password flag",
			inputMasterPass: "",
			masterPassValid: false,
			expectedOutput:  "❌ Ошибка: укажите пароль (--masterPass или -m)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			cmd := &cobra.Command{}
			cmd.Flags().StringP("masterPass", "m", tt.inputMasterPass, "Мастер пароль")

			if tt.inputMasterPass != "" && tt.masterPassValid {
				mockUserService.EXPECT().Master(tt.inputMasterPass).Return(nil)
			} else if tt.inputMasterPass != "" && !tt.masterPassValid {
				mockUserService.EXPECT().Master(tt.inputMasterPass).Return(fmt.Errorf("invalid master password"))
			}

			userHandler.masterRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}

func TestUserHandler_LoginRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := service.NewMockUserService(ctrl)
	userHandler := NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		login          string
		password       string
		mockResponse   string
		mockError      error
		expectedOutput string
	}{
		{
			name:           "successful login",
			login:          "testuser",
			password:       "password123",
			mockResponse:   "valid_token",
			mockError:      nil,
			expectedOutput: "✅ Авторизация прошла успешно",
		},
		{
			name:           "login with empty response token",
			login:          "testuser",
			password:       "password1234",
			mockResponse:   "",
			mockError:      errors.New("Не удалось получить токен"),
			expectedOutput: "❌ Ошибка авторизации: Не удалось получить токен",
		},
		{
			name:           "login with service error",
			login:          "testuser",
			password:       "password1235",
			mockResponse:   "",
			mockError:      errors.New("invalid credentials"),
			expectedOutput: "❌ Ошибка авторизации: invalid credentials",
		},
		{
			name:           "missing login",
			login:          "",
			password:       "password123",
			mockResponse:   "",
			mockError:      nil,
			expectedOutput: "❌ Ошибка: укажите логин (--login или -l)",
		},
		{
			name:           "missing password",
			login:          "testuser",
			password:       "",
			mockResponse:   "",
			mockError:      nil,
			expectedOutput: "❌ Ошибка: укажите пароль (--password или -p)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			cmd := &cobra.Command{}
			cmd.Flags().StringP("login", "l", tt.login, "Логин")
			cmd.Flags().StringP("password", "p", tt.password, "Пароль")

			if tt.login != "" && tt.password != "" {
				mockUserService.EXPECT().Login(tt.login, tt.password).Return(tt.mockResponse, tt.mockError).AnyTimes()
			}

			userHandler.loginRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}

func TestUserHandler_RegistrationRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := service.NewMockUserService(ctrl)
	userHandler := NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		login          string
		password       string
		masterPass     string
		mockResponse   string
		mockError      error
		expectedOutput string
	}{
		{
			name:           "successful registration",
			login:          "testuser",
			password:       "password123",
			masterPass:     "masterpass123",
			mockResponse:   "valid_token",
			mockError:      nil,
			expectedOutput: "✅ Регистрация прошла успешно",
		},
		{
			name:           "registration with empty response token",
			login:          "testuser2",
			password:       "password123",
			masterPass:     "masterpass123",
			mockResponse:   "",
			mockError:      nil,
			expectedOutput: "❌ Ошибка регистрации: Не удалось получить токен",
		},
		{
			name:           "registration with service error",
			login:          "testuser3",
			password:       "password123",
			masterPass:     "masterpass123",
			mockResponse:   "",
			mockError:      errors.New("user already exists"),
			expectedOutput: "❌ Ошибка регистрации: user already exists",
		},
		{
			name:           "missing login",
			login:          "",
			password:       "password123",
			masterPass:     "masterpass123",
			mockResponse:   "",
			mockError:      nil,
			expectedOutput: "❌ Ошибка: укажите логин (--login или -l)",
		},
		{
			name:           "missing password",
			login:          "testuser6",
			password:       "",
			masterPass:     "masterpass123",
			mockResponse:   "",
			mockError:      nil,
			expectedOutput: "❌ Ошибка: укажите пароль (--password или -p)",
		},
		{
			name:           "missing master password",
			login:          "testuser7",
			password:       "password123",
			masterPass:     "",
			mockResponse:   "",
			mockError:      nil,
			expectedOutput: "❌ Ошибка: укажите пароль (--masterPass или -m)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			cmd := &cobra.Command{}
			cmd.Flags().StringP("login", "l", tt.login, "Логин")
			cmd.Flags().StringP("password", "p", tt.password, "Пароль")
			cmd.Flags().StringP("masterPass", "m", tt.masterPass, "Мастер пароль")

			if tt.login != "" && tt.password != "" && tt.masterPass != "" {
				mockUserService.EXPECT().Register(tt.login, tt.password, tt.masterPass).Return(tt.mockResponse, tt.mockError).AnyTimes()
			}

			userHandler.registrationRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}
