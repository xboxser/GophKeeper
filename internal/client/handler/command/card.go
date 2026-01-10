package command

import (
	"fmt"
	"gophkeeper/internal/client/service"
	"gophkeeper/internal/model"

	"github.com/spf13/cobra"
)

type CardHandler struct {
	CardService service.CardService
}

func NewCardHandler(cardService service.CardService) *CardHandler {
	return &CardHandler{
		CardService: cardService,
	}

}
func (ch *CardHandler) GetCommands() []*cobra.Command {
	addCardCmd := &cobra.Command{
		Use:     "addCard",
		Short:   "Добавить новую банковскую карту",
		Example: `  todo addCard --title "Название карты" --number "1234 1234 1234 1234" --expiry "12/23" --card_holder "Имя и фамилия" --cvv "123"`,
		Run:     ch.addCardRun,
	}
	addCardCmd.Flags().StringP("title", "t", "", "Произвольное название карты (обязательное)")
	addCardCmd.MarkFlagRequired("title")
	addCardCmd.Flags().StringP("number", "n", "", "Номер карты (обязательное)")
	addCardCmd.MarkFlagRequired("number")
	addCardCmd.Flags().StringP("expiry", "e", "", "Срок действия карты (обязательное)")
	addCardCmd.MarkFlagRequired("expiry")
	addCardCmd.Flags().StringP("cvv", "c", "", "cvv (обязательное)")
	addCardCmd.MarkFlagRequired("cvv")
	addCardCmd.Flags().StringP("card_holder", "o", "", "Имя держателя карты (обязательное)")
	addCardCmd.MarkFlagRequired("card_holder")

	getCardCmd := &cobra.Command{
		Use:     "getCard",
		Short:   "Получить список банковских карт",
		Example: ``,
		Run:     ch.getCardRun,
	}

	return []*cobra.Command{
		addCardCmd,
		getCardCmd,
	}
}

func (ch *CardHandler) getCardRun(cmd *cobra.Command, args []string) {
	cards, err := ch.CardService.GetCards()

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	if len(cards) == 0 {
		fmt.Println("📦 Нет сохраненных банковских карт")
		return
	}

	fmt.Printf("🔐 Найдено %d записей банковских карт:\n\n", len(cards))

	for i, card := range cards {
		fmt.Printf("📋 Карта #%d\n", i+1)
		fmt.Printf("Title: %s\n", card.Title)
		fmt.Printf("Number: %s\n", card.Number)
		fmt.Printf("Expiry: %s\n", card.Expiry)
		fmt.Printf("Card Holder: %s\n", card.CardHolder)
		fmt.Printf("CVV: %s\n", card.CVV)
		fmt.Println()
	}
}

func (ch *CardHandler) addCardRun(cmd *cobra.Command, args []string) {
	card := model.Card{}

	str, err := cmd.Flags().GetString("title")
	if err != nil || str == "" {
		fmt.Println("❌ Ошибка: укажите произвольное название карты (--title или -t)")
		return
	}
	card.Title = str

	str, err = cmd.Flags().GetString("number")
	if err != nil || str == "" {
		fmt.Println("❌ Ошибка: укажите номер карты (--number или -n)")
		return
	}
	card.Number = str

	str, err = cmd.Flags().GetString("expiry")
	if err != nil || str == "" {
		fmt.Println("❌ Ошибка: укажите срок действия карты (--expiry или -e)")
		return
	}
	card.Expiry = str

	str, err = cmd.Flags().GetString("cvv")
	if err != nil || str == "" {
		fmt.Println("❌ Ошибка: укажите cvv (--cvv или -c)")
		return
	}
	card.CVV = str

	str, err = cmd.Flags().GetString("card_holder")
	if err != nil || str == "" {
		fmt.Println("❌ Ошибка: укажите card_holder (--card_holder или -o)")
		return
	}
	card.CardHolder = str
	err = ch.CardService.AddCard(card)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	fmt.Printf("✅ Запись успешно добавлена\n")

}
