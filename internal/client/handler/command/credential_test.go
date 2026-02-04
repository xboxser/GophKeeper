package command

import (
	"fmt"
	"gophkeeper/internal/model"
	service "gophkeeper/mocks/client/services"
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCredentialHandler_GetCommands(t *testing.T) {
	mockCredentialService := &service.MockCredentialService{}
	credentialHandler := NewCredentialHandler(mockCredentialService, nil)

	commands := credentialHandler.GetCommands()

	assert.Len(t, commands, 4)

	commandUses := make([]string, 0, 4)
	for _, cmd := range commands {
		commandUses = append(commandUses, cmd.Use)
	}

	assert.Contains(t, commandUses, "credential-add")
	assert.Contains(t, commandUses, "credential-get")
	assert.Contains(t, commandUses, "credential-del")
	assert.Contains(t, commandUses, "credential-update")
}

func TestCredentialHandler_AddCredentialRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredentialService := service.NewMockCredentialService(ctrl)
	mockUserMasterService := service.NewMockUserMasterService(ctrl)
	credentialHandler := NewCredentialHandler(mockCredentialService, mockUserMasterService)

	tests := []struct {
		name           string
		login          string
		password       string
		masterPass     string
		masterValid    bool
		serviceError   error
		expectedOutput string
	}{
		{
			name:           "successful credential addition",
			login:          "testuser",
			password:       "password123",
			masterPass:     "masterpass123",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "✅ Запись успешно добавлена",
		},
		{
			name:           "invalid master password",
			login:          "testuser",
			password:       "password123",
			masterPass:     "wrongmasterpass",
			masterValid:    false,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка проверки мастер пароля:",
		},
		{
			name:           "service error during credential addition",
			login:          "testuser",
			password:       "password123",
			masterPass:     "masterpass123",
			masterValid:    true,
			serviceError:   fmt.Errorf("database error"),
			expectedOutput: "❌ Ошибка: database error",
		},
		{
			name:           "missing login",
			login:          "",
			password:       "password123",
			masterPass:     "masterpass123",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка: укажите логин (--login или -l)",
		},
		{
			name:           "missing password",
			login:          "testuser",
			password:       "",
			masterPass:     "masterpass123",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка: укажите пароль (--password или -p)",
		},
		{
			name:           "missing master password",
			login:          "testuser",
			password:       "password123",
			masterPass:     "",
			masterValid:    true,
			serviceError:   nil,
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
				if tt.masterValid {
					mockUserMasterService.EXPECT().Master(tt.masterPass).Return(nil)
					if tt.serviceError == nil {
						mockCredentialService.EXPECT().AddCredential(tt.login, tt.password, tt.masterPass).Return(nil)
					} else {
						mockCredentialService.EXPECT().AddCredential(tt.login, tt.password, tt.masterPass).Return(tt.serviceError)
					}
				} else {
					mockUserMasterService.EXPECT().Master(tt.masterPass).Return(fmt.Errorf("invalid master password"))
				}
			}

			credentialHandler.addCredentialRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}

func TestCredentialHandler_DeleteCredentialRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredentialService := service.NewMockCredentialService(ctrl)
	mockUserMasterService := service.NewMockUserMasterService(ctrl)
	credentialHandler := NewCredentialHandler(mockCredentialService, mockUserMasterService)

	tests := []struct {
		name           string
		masterPass     string
		recordID       string
		masterValid    bool
		serviceError   error
		expectedOutput string
	}{
		{
			name:           "successful credential deletion",
			masterPass:     "masterpass123",
			recordID:       "1",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "✅ Запись успешно удалена",
		},
		{
			name:           "invalid master password",
			masterPass:     "wrongmasterpass",
			recordID:       "12",
			masterValid:    false,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка проверки мастер пароля:",
		},
		{
			name:           "service error during credential deletion",
			masterPass:     "masterpass123",
			recordID:       "13",
			masterValid:    true,
			serviceError:   fmt.Errorf("record not found"),
			expectedOutput: "❌ Ошибка: record not found",
		},
		{
			name:           "missing master password",
			masterPass:     "",
			recordID:       "14",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка: укажите пароль (--masterPass или -m)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			cmd := &cobra.Command{}
			cmd.Flags().StringP("masterPass", "m", tt.masterPass, "Мастер пароль")
			cmd.Flags().StringP("id", "i", tt.recordID, "Номер записи")

			if tt.masterPass != "" && tt.recordID != "" {
				if tt.masterValid {
					mockUserMasterService.EXPECT().Master(tt.masterPass).Return(nil)
					if tt.serviceError == nil {
						mockCredentialService.EXPECT().DeleteCredential(tt.recordID).Return(nil)
					} else {
						mockCredentialService.EXPECT().DeleteCredential(tt.recordID).Return(tt.serviceError)
					}
				} else {
					mockUserMasterService.EXPECT().Master(tt.masterPass).Return(fmt.Errorf("invalid master password"))
				}
			}

			credentialHandler.deleteCredentialRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}

func TestCredentialHandler_GetCredentialRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredentialService := service.NewMockCredentialService(ctrl)
	mockUserMasterService := service.NewMockUserMasterService(ctrl) // Even though not used directly, we need it for the constructor
	credentialHandler := NewCredentialHandler(mockCredentialService, mockUserMasterService)

	tests := []struct {
		name            string
		masterPass      string
		serviceResponse []model.Credential
		serviceError    error
		expectedOutput  string
	}{
		{
			name:       "successful retrieval with credentials",
			masterPass: "masterpass123",
			serviceResponse: []model.Credential{
				{ID: 1, Login: "testuser1", Password: "password1"},
				{ID: 2, Login: "testuser2", Password: "password2"},
			},
			serviceError:   nil,
			expectedOutput: "🔐 Найдено 2 учетных записей:",
		},
		{
			name:            "no credentials found",
			masterPass:      "masterpass123",
			serviceResponse: []model.Credential{},
			serviceError:    nil,
			expectedOutput:  "📦 Нет сохраненных учетных данных",
		},
		{
			name:            "service error during credential retrieval",
			masterPass:      "masterpass123",
			serviceResponse: nil,
			serviceError:    fmt.Errorf("database connection failed"),
			expectedOutput:  "❌ Ошибка: database connection failed",
		},
		{
			name:            "missing master password",
			masterPass:      "",
			serviceResponse: nil,
			serviceError:    nil,
			expectedOutput:  "❌ Ошибка: укажите пароль (--masterPass или -m)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			cmd := &cobra.Command{}
			cmd.Flags().StringP("masterPass", "m", tt.masterPass, "Мастер пароль")

			if tt.masterPass != "" {
				if tt.serviceError != nil {
					mockCredentialService.EXPECT().GetCredentials(tt.masterPass).Return(nil, tt.serviceError)
				} else {
					mockCredentialService.EXPECT().GetCredentials(tt.masterPass).Return(tt.serviceResponse, nil)
				}
			}

			credentialHandler.getCredentialRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}

func TestCredentialHandler_UpdateCredentialRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredentialService := service.NewMockCredentialService(ctrl)
	mockUserMasterService := service.NewMockUserMasterService(ctrl)
	credentialHandler := NewCredentialHandler(mockCredentialService, mockUserMasterService)

	tests := []struct {
		name           string
		login          string
		password       string
		masterPass     string
		id             string
		masterValid    bool
		serviceError   error
		expectedOutput string
	}{
		{
			name:           "successful credential update",
			login:          "testuser",
			password:       "newpassword123",
			masterPass:     "masterpass123",
			id:             "1",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "✅ Запись успешно обновлена",
		},
		{
			name:           "invalid master password",
			login:          "testuser2",
			password:       "newpassword123",
			masterPass:     "wrongmasterpass",
			id:             "1",
			masterValid:    false,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка проверки мастер пароля:",
		},
		{
			name:           "service error during credential update",
			login:          "testuser3",
			password:       "newpassword123",
			masterPass:     "masterpass123",
			id:             "1",
			masterValid:    true,
			serviceError:   fmt.Errorf("record not found"),
			expectedOutput: "❌ Ошибка: record not found",
		},
		{
			name:           "missing login",
			login:          "",
			password:       "newpassword123",
			masterPass:     "masterpass123",
			id:             "1",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка: укажите логин (--login или -l)",
		},
		{
			name:           "missing password",
			login:          "testuser4",
			password:       "",
			masterPass:     "masterpass123",
			id:             "1",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка: укажите пароль (--password или -p)",
		},
		{
			name:           "missing master password",
			login:          "testuser5",
			password:       "newpassword123",
			masterPass:     "",
			id:             "1",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка: укажите пароль (--masterPass или -m)",
		},
		{
			name:           "missing record ID",
			login:          "testuser6",
			password:       "newpassword123",
			masterPass:     "masterpass123",
			id:             "",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка: номер записи (--id или -i)",
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
			cmd.Flags().StringP("id", "i", tt.id, "Номер записи")

			if tt.login != "" && tt.password != "" && tt.masterPass != "" && tt.id != "" {
				if tt.masterValid {
					mockUserMasterService.EXPECT().Master(tt.masterPass).Return(nil)
					if tt.serviceError == nil {
						mockCredentialService.EXPECT().UpdateCredential(tt.login, tt.password, tt.masterPass, tt.id).Return(nil)
					} else {
						mockCredentialService.EXPECT().UpdateCredential(tt.login, tt.password, tt.masterPass, tt.id).Return(tt.serviceError)
					}
				} else {
					mockUserMasterService.EXPECT().Master(tt.masterPass).Return(fmt.Errorf("invalid master password"))
				}
			}

			credentialHandler.updateCredentialRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}
