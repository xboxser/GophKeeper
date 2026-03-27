package service

import (
	"encoding/json"
	"gophkeeper/internal/model"
	"gophkeeper/mocks/client/service"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCardSetMasterPass(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	cardService := NewCardService(mockSender, mockEncryption)

	masterPass := "test_master_password"

	mockEncryption.EXPECT().SetMasterPass(masterPass)

	cardService.SetMasterPass(masterPass)
}

func TestDeleteCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	service := NewCardService(mockSender, mockEncryption)
	service.InitToken(mockToken)

	cardID := "123"

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")
	mockSender.EXPECT().SendDelete(gomock.Any(), "/api/card/"+cardID).Return(
		[]byte("success"),
		&http.Response{StatusCode: http.StatusNoContent},
		nil,
	)

	err := service.DeleteCard(cardID)

	require.NoError(t, err)
}

func TestGetCards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	service := NewCardService(mockSender, mockEncryption)
	service.InitToken(mockToken)

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")

	expectedCardsAPI := []model.CardAPI{
		{
			Title:      "Test Card",
			Number:     []byte("encrypted_number"),
			CVV:        []byte("encrypted_cvv"),
			Expiry:     []byte("encrypted_expiry"),
			CardHolder: []byte("encrypted_holder"),
		},
	}
	jsonData, _ := json.Marshal(expectedCardsAPI)

	mockSender.EXPECT().SendGet(gomock.Any(), "/api/card").Return(
		jsonData,
		&http.Response{StatusCode: http.StatusOK},
		nil,
	)

	mockEncryption.EXPECT().Decrypt([]byte("encrypted_number")).Return("1234567890123456", nil)
	mockEncryption.EXPECT().Decrypt([]byte("encrypted_cvv")).Return("123", nil)
	mockEncryption.EXPECT().Decrypt([]byte("encrypted_expiry")).Return("12/25", nil)
	mockEncryption.EXPECT().Decrypt([]byte("encrypted_holder")).Return("John Doe", nil)

	result, err := service.GetCards()

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "Test Card", result[0].Title)
}

func TestAddCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	cardService := NewCardService(mockSender, mockEncryption)
	cardService.InitToken(mockToken)

	testCard := model.Card{
		ID:         1,
		Title:      "Test Card",
		Number:     "1234567890123456",
		CVV:        "123",
		Expiry:     "12/25",
		CardHolder: "John Doe",
	}

	expectedCardAPI := model.CardAPI{
		IncID:      1,
		Title:      "Test Card",
		Last4:      "3456", // Last 4 digits of the card number
		Number:     []byte("encrypted_number"),
		CVV:        []byte("encrypted_cvv"),
		Expiry:     []byte("encrypted_expiry"),
		CardHolder: []byte("encrypted_holder"),
	}

	jsonData, _ := json.Marshal(expectedCardAPI)

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")

	mockEncryption.EXPECT().Encrypt(testCard.Number).Return([]byte("encrypted_number"), nil)
	mockEncryption.EXPECT().Encrypt(testCard.CVV).Return([]byte("encrypted_cvv"), nil)
	mockEncryption.EXPECT().Encrypt(testCard.Expiry).Return([]byte("encrypted_expiry"), nil)
	mockEncryption.EXPECT().Encrypt(testCard.CardHolder).Return([]byte("encrypted_holder"), nil)

	mockSender.EXPECT().SendPost(gomock.Any(), "/api/card", jsonData).Return(
		[]byte("success"),
		&http.Response{StatusCode: http.StatusOK},
		nil,
	)
	err := cardService.AddCard(testCard)

	require.NoError(t, err)
}

func TestUpdateCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSender := service.NewMockSenderService(ctrl)
	mockToken := service.NewMockTokenService(ctrl)
	mockEncryption := service.NewMockEncryptionService(ctrl)

	cardService := NewCardService(mockSender, mockEncryption)
	cardService.InitToken(mockToken)

	testCard := model.Card{
		ID:         1,
		Title:      "Updated Card",
		Number:     "1234567890123456",
		CVV:        "123",
		Expiry:     "12/25",
		CardHolder: "John Doe",
	}

	expectedCardAPI := model.CardAPI{
		IncID:      1,
		Title:      "Updated Card",
		Last4:      "3456",
		Number:     []byte("encrypted_number"),
		CVV:        []byte("encrypted_cvv"),
		Expiry:     []byte("encrypted_expiry"),
		CardHolder: []byte("encrypted_holder"),
	}

	jsonData, _ := json.Marshal(expectedCardAPI)

	mockToken.EXPECT().GetToken().Return("token", nil)
	mockSender.EXPECT().SetToken("token")

	mockEncryption.EXPECT().Encrypt(testCard.Number).Return([]byte("encrypted_number"), nil)
	mockEncryption.EXPECT().Encrypt(testCard.CVV).Return([]byte("encrypted_cvv"), nil)
	mockEncryption.EXPECT().Encrypt(testCard.Expiry).Return([]byte("encrypted_expiry"), nil)
	mockEncryption.EXPECT().Encrypt(testCard.CardHolder).Return([]byte("encrypted_holder"), nil)

	mockSender.EXPECT().SendPut(gomock.Any(), "/api/card", jsonData).Return(
		[]byte("success"),
		&http.Response{StatusCode: http.StatusOK},
		nil,
	)

	err := cardService.UpdateCard(testCard)

	require.NoError(t, err)
}
