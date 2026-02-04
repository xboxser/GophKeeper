package command

import (
	"fmt"
	"gophkeeper/internal/model"
	"gophkeeper/mocks/client/service"
	services "gophkeeper/mocks/client/services"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCardHandler_GetCommands(t *testing.T) {
	mockCardService := &service.MockCardService{}
	cardHandler := NewCardHandler(mockCardService, nil)

	commands := cardHandler.GetCommands()

	assert.Len(t, commands, 4)

	commandUses := make([]string, 0, 4)
	for _, cmd := range commands {
		commandUses = append(commandUses, cmd.Use)
	}

	assert.Contains(t, commandUses, "card-add")
	assert.Contains(t, commandUses, "card-get")
	assert.Contains(t, commandUses, "card-delete")
	assert.Contains(t, commandUses, "card-update")
}

func TestCardHandler_DeleteCardRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := service.NewMockCardService(ctrl)
	mockUserMasterService := services.NewMockUserMasterService(ctrl)
	cardHandler := NewCardHandler(mockCardService, mockUserMasterService)

	tests := []struct {
		name           string
		masterPass     string
		cardID         string
		masterValid    bool
		serviceError   error
		expectedOutput string
	}{
		{
			name:           "successful card deletion",
			masterPass:     "masterpass123",
			cardID:         "1",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "✅ Карта успешно удалена",
		},
		{
			name:           "invalid master password",
			masterPass:     "wrongmasterpass",
			cardID:         "1",
			masterValid:    false,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка проверки мастер пароля:",
		},
		{
			name:           "service error during card deletion",
			masterPass:     "masterpass123",
			cardID:         "1",
			masterValid:    true,
			serviceError:   fmt.Errorf("card not found"),
			expectedOutput: "❌ Ошибка: card not found",
		},
		{
			name:           "missing master password",
			masterPass:     "",
			cardID:         "1",
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
			cmd.Flags().StringP("id", "i", tt.cardID, "ID карты")

			// Only expect SetMasterPass and Master validation if both required fields are provided
			if tt.masterPass != "" && tt.cardID != "" {
				mockCardService.EXPECT().SetMasterPass(tt.masterPass)
				if tt.masterValid {
					if tt.serviceError == nil {
						mockUserMasterService.EXPECT().Master(tt.masterPass).Return(nil)
						mockCardService.EXPECT().DeleteCard(tt.cardID).Return(nil)
					} else {
						mockUserMasterService.EXPECT().Master(tt.masterPass).Return(nil)
						mockCardService.EXPECT().DeleteCard(tt.cardID).Return(tt.serviceError)
					}
				} else {
					mockUserMasterService.EXPECT().Master(tt.masterPass).Return(fmt.Errorf("invalid master password"))
				}
			}

			cardHandler.deleteCardRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}

func TestCardHandler_GetCardRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := service.NewMockCardService(ctrl)
	mockUserMasterService := services.NewMockUserMasterService(ctrl) // Even though not used directly, we need it for the constructor
	cardHandler := NewCardHandler(mockCardService, mockUserMasterService)

	tests := []struct {
		name            string
		masterPass      string
		serviceResponse []model.Card
		serviceError    error
		expectedOutput  string
	}{
		{
			name:       "successful retrieval with cards",
			masterPass: "masterpass123",
			serviceResponse: []model.Card{
				{ID: 1, Title: "My Visa Card", Number: "1234 5678 9012 3456", Expiry: "12/25", CardHolder: "John Doe", CVV: "123"},
				{ID: 2, Title: "Business Mastercard", Number: "9876 5432 1098 7654", Expiry: "08/24", CardHolder: "Jane Smith", CVV: "456"},
			},
			serviceError:   nil,
			expectedOutput: "🔐 Найдено 2 записей банковских карт:",
		},
		{
			name:            "no cards found",
			masterPass:      "masterpass123",
			serviceResponse: []model.Card{},
			serviceError:    nil,
			expectedOutput:  "📦 Нет сохраненных банковских карт",
		},
		{
			name:            "service error during card retrieval",
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
				mockCardService.EXPECT().SetMasterPass(tt.masterPass)
				if tt.serviceError != nil {
					mockCardService.EXPECT().GetCards().Return(nil, tt.serviceError)
				} else {
					mockCardService.EXPECT().GetCards().Return(tt.serviceResponse, nil)
				}
			}

			cardHandler.getCardRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}

func TestCardHandler_AddCardRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardService := service.NewMockCardService(ctrl)
	mockUserMasterService := services.NewMockUserMasterService(ctrl)
	cardHandler := NewCardHandler(mockCardService, mockUserMasterService)

	tests := []struct {
		name           string
		title          string
		number         string
		expiry         string
		cvv            string
		cardHolder     string
		masterPass     string
		masterValid    bool
		serviceError   error
		expectedOutput string
	}{
		{
			name:           "successful card addition",
			title:          "My Visa Card",
			number:         "1234 5678 9012 3456",
			expiry:         "12/25",
			cvv:            "123",
			cardHolder:     "John Doe",
			masterPass:     "masterpass123",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "✅ Запись успешно добавлена",
		},
		{
			name:           "invalid master password",
			title:          "My Visa Card",
			number:         "1234 5678 9012 3456",
			expiry:         "12/25",
			cvv:            "123",
			cardHolder:     "John Doe",
			masterPass:     "wrongmasterpass",
			masterValid:    false,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка проверки мастер пароля:",
		},
		{
			name:           "service error during card addition",
			title:          "My Visa Card",
			number:         "1234 5678 9012 3456",
			expiry:         "12/25",
			cvv:            "123",
			cardHolder:     "John Doe",
			masterPass:     "masterpass123",
			masterValid:    true,
			serviceError:   fmt.Errorf("database error"),
			expectedOutput: "❌ Ошибка: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			cmd := &cobra.Command{}
			cmd.Flags().StringP("title", "t", tt.title, "Произвольное название карты")
			cmd.Flags().StringP("number", "n", tt.number, "Номер карты")
			cmd.Flags().StringP("expiry", "e", tt.expiry, "Срок действия карты")
			cmd.Flags().StringP("cvv", "c", tt.cvv, "CVV")
			cmd.Flags().StringP("card_holder", "o", tt.cardHolder, "Имя держателя карты")
			cmd.Flags().StringP("masterPass", "m", tt.masterPass, "Мастер пароль")

			// Only expect SetMasterPass and Master validation if master password is provided
			if tt.masterPass != "" {
				mockCardService.EXPECT().SetMasterPass(tt.masterPass)
				if tt.masterValid {
					if tt.title != "" && tt.number != "" && tt.expiry != "" && tt.cvv != "" && tt.cardHolder != "" {
						mockUserMasterService.EXPECT().Master(tt.masterPass).Return(nil)
						expectedCard := model.Card{
							Title:      tt.title,
							Number:     tt.number,
							Expiry:     tt.expiry,
							CVV:        tt.cvv,
							CardHolder: tt.cardHolder,
						}
						if tt.serviceError == nil {
							mockCardService.EXPECT().AddCard(expectedCard).Return(nil)
						} else {
							mockCardService.EXPECT().AddCard(expectedCard).Return(tt.serviceError)
						}
					}
				} else {
					mockUserMasterService.EXPECT().Master(tt.masterPass).Return(fmt.Errorf("invalid master password"))
				}
			}

			cardHandler.addCardRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}

func TestCardHandler_UpdateCardRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		title          string
		id             string
		number         string
		expiry         string
		cvv            string
		cardHolder     string
		masterPass     string
		masterValid    bool
		serviceError   error
		expectedOutput string
	}{
		{
			name:           "successful card update",
			title:          "Updated Visa Card",
			id:             "1",
			number:         "1234 5678 9012 3456",
			expiry:         "12/25",
			cvv:            "123",
			cardHolder:     "John Doe",
			masterPass:     "masterpass123",
			masterValid:    true,
			serviceError:   nil,
			expectedOutput: "✅ Запись успешно обновлена",
		},
		{
			name:           "invalid master password",
			title:          "Updated Visa Card",
			id:             "1",
			number:         "1234 5678 9012 3456",
			expiry:         "12/25",
			cvv:            "123",
			cardHolder:     "John Doe",
			masterPass:     "wrongmasterpass",
			masterValid:    false,
			serviceError:   nil,
			expectedOutput: "❌ Ошибка проверки мастер пароля:",
		},
		{
			name:           "service error during card update",
			title:          "Updated Visa Card",
			id:             "1",
			number:         "1234 5678 9012 3456",
			expiry:         "12/25",
			cvv:            "123",
			cardHolder:     "John Doe",
			masterPass:     "masterpass123",
			masterValid:    true,
			serviceError:   fmt.Errorf("card not found"),
			expectedOutput: "❌ Ошибка: card not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new controller for each subtest to avoid conflicts
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCardService := service.NewMockCardService(ctrl)
			mockUserMasterService := services.NewMockUserMasterService(ctrl)
			cardHandler := NewCardHandler(mockCardService, mockUserMasterService)

			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			cmd := &cobra.Command{}
			cmd.Flags().StringP("title", "t", tt.title, "Произвольное название карты")
			cmd.Flags().StringP("id", "i", tt.id, "ID карты")
			cmd.Flags().StringP("number", "n", tt.number, "Номер карты")
			cmd.Flags().StringP("expiry", "e", tt.expiry, "Срок действия карты")
			cmd.Flags().StringP("cvv", "c", tt.cvv, "CVV")
			cmd.Flags().StringP("card_holder", "o", tt.cardHolder, "Имя держателя карты")
			cmd.Flags().StringP("masterPass", "m", tt.masterPass, "Мастер пароль")

			// Setup mock expectations only when needed
			if tt.masterPass != "" {
				mockCardService.EXPECT().SetMasterPass(tt.masterPass).AnyTimes()

				// Check if we have all required fields except for validation failures
				hasAllRequiredFields := tt.title != "" && tt.id != "" && tt.number != "" &&
					tt.expiry != "" && tt.cvv != "" && tt.cardHolder != ""

				if hasAllRequiredFields {
					// Try to parse the ID to see if it's valid
					_, idParseErr := strconv.Atoi(tt.id)

					if idParseErr != nil {
						// ID parsing will fail, so no more expectations needed after this point
					} else if tt.masterValid {
						// Expect master validation to succeed
						mockUserMasterService.EXPECT().Master(tt.masterPass).Return(nil).AnyTimes()

						// Only expect UpdateCard if validation passes
						if tt.serviceError == nil {
							expectedCard := model.Card{
								ID:         1, // This assumes the first successful test case with ID "1"
								Title:      tt.title,
								Number:     tt.number,
								Expiry:     tt.expiry,
								CVV:        tt.cvv,
								CardHolder: tt.cardHolder,
							}
							// Convert string ID to int for the expected card
							if parsedId, err := strconv.Atoi(tt.id); err == nil {
								expectedCard.ID = parsedId
								mockCardService.EXPECT().UpdateCard(expectedCard).Return(nil).AnyTimes()
							}
						} else {
							if parsedId, err := strconv.Atoi(tt.id); err == nil {
								expectedCard := model.Card{
									ID:         parsedId,
									Title:      tt.title,
									Number:     tt.number,
									Expiry:     tt.expiry,
									CVV:        tt.cvv,
									CardHolder: tt.cardHolder,
								}
								mockCardService.EXPECT().UpdateCard(expectedCard).Return(tt.serviceError).AnyTimes()
							}
						}
					} else {
						// Expect master validation to fail
						mockUserMasterService.EXPECT().Master(tt.masterPass).Return(fmt.Errorf("invalid master password")).AnyTimes()
					}
				}
			}

			cardHandler.updateCardRun(cmd, []string{})

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			outputStr := string(out)
			assert.Contains(t, outputStr, tt.expectedOutput)
		})
	}
}
