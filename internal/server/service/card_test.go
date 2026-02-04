package service

import (
	"context"
	"errors"
	"gophkeeper/internal/model"
	"gophkeeper/mocks/server/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardRepo := repository.NewMockCardRepository(ctrl)
	cardService := NewCardService(mockCardRepo)

	userID := 123

	t.Run("success get cards", func(t *testing.T) {
		expectedCards := []model.CardAPI{
			{
				ID:         1,
				Title:      "Card 1",
				Number:     []byte("1234567890123456"),
				Expiry:     []byte("12/24"),
				CVV:        []byte("123"),
				CardHolder: []byte("Ivan"),
				Last4:      "3456",
				IncID:      1,
			},
		}

		mockCardRepo.EXPECT().GetCards(gomock.Any(), userID).Return(expectedCards, nil)

		// Вызываем тестируемый метод
		result, err := cardService.GetCards(context.Background(), userID)

		// Проверяем результат
		require.NoError(t, err)
		assert.Equal(t, expectedCards, result)
	})

	t.Run("error getting cards", func(t *testing.T) {
		expectedError := errors.New("database error")

		mockCardRepo.EXPECT().GetCards(gomock.Any(), userID).Return(nil, expectedError)

		result, err := cardService.GetCards(context.Background(), userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})

}

func TestAddCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardRepo := repository.NewMockCardRepository(ctrl)
	cardService := NewCardService(mockCardRepo)

	userID := 123
	card := &model.CardAPI{
		ID:         1,
		Title:      "Test Card",
		Number:     []byte("1234567890123456"),
		Expiry:     []byte("12/25"),
		CVV:        []byte("123"),
		CardHolder: []byte("John Doe"),
		Last4:      "3456",
		IncID:      1,
	}

	t.Run("successful add card", func(t *testing.T) {
		mockCardRepo.EXPECT().AddCard(gomock.Any(), *card, userID).Return(nil)

		err := cardService.AddCard(context.Background(), card, userID)

		require.NoError(t, err)
	})

	t.Run("error adding card", func(t *testing.T) {
		expectedError := errors.New("database error")

		mockCardRepo.EXPECT().AddCard(gomock.Any(), *card, userID).Return(expectedError)

		err := cardService.AddCard(context.Background(), card, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestDeleteCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardRepo := repository.NewMockCardRepository(ctrl)
	cardService := NewCardService(mockCardRepo)

	userID := 123
	cardID := "1"

	t.Run("successful delete card", func(t *testing.T) {
		// Подготовим ожидаемые данные
		expectedCard := model.CardAPI{
			ID:         1,
			Title:      "Card 1",
			Number:     []byte("1234567890123456"),
			Expiry:     []byte("12/24"),
			CVV:        []byte("123"),
			CardHolder: []byte("Ivan"),
			Last4:      "3456",
			IncID:      1,
		}

		mockCardRepo.EXPECT().GetCard(gomock.Any(), cardID, userID).Return(expectedCard, nil)
		mockCardRepo.EXPECT().DeleteCard(gomock.Any(), cardID, userID).Return(nil)

		err := cardService.DeleteCard(context.Background(), cardID, userID)

		require.NoError(t, err)
	})

	t.Run("card not found", func(t *testing.T) {
		emptyCard := model.CardAPI{}

		mockCardRepo.EXPECT().GetCard(gomock.Any(), cardID, userID).Return(emptyCard, nil)

		err := cardService.DeleteCard(context.Background(), cardID, userID)

		assert.Error(t, err)
		assert.Equal(t, model.ErrCardNotFound, err)
	})

	t.Run("error getting card", func(t *testing.T) {
		expectedError := errors.New("database error")

		mockCardRepo.EXPECT().GetCard(gomock.Any(), cardID, userID).Return(model.CardAPI{}, expectedError)

		err := cardService.DeleteCard(context.Background(), cardID, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestUpdateCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCardRepo := repository.NewMockCardRepository(ctrl)
	cardService := NewCardService(mockCardRepo)

	userID := 123
	cardID := "1"
	// Подготовим ожидаемые данные
	expectedCard := model.CardAPI{
		ID:         1,
		Title:      "Card 1",
		Number:     []byte("1234567890123456"),
		Expiry:     []byte("12/24"),
		CVV:        []byte("123"),
		CardHolder: []byte("Ivan"),
		Last4:      "3456",
		IncID:      1,
	}
	t.Run("successful update card", func(t *testing.T) {

		mockCardRepo.EXPECT().GetCard(gomock.Any(), cardID, userID).Return(expectedCard, nil)
		mockCardRepo.EXPECT().UpdateCard(gomock.Any(), expectedCard, userID).Return(nil)

		err := cardService.UpdateCard(context.Background(), &expectedCard, userID)

		require.NoError(t, err)
	})

	t.Run("card not found", func(t *testing.T) {
		emptyCard := model.CardAPI{}

		mockCardRepo.EXPECT().GetCard(gomock.Any(), cardID, userID).Return(emptyCard, nil)

		err := cardService.UpdateCard(context.Background(), &expectedCard, userID)

		assert.Error(t, err)
		assert.Equal(t, model.ErrCardNotFound, err)
	})

	t.Run("error getting card", func(t *testing.T) {
		expectedError := errors.New("database error")

		mockCardRepo.EXPECT().GetCard(gomock.Any(), cardID, userID).Return(model.CardAPI{}, expectedError)

		err := cardService.UpdateCard(context.Background(), &expectedCard, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}
