package command

import (
	"fmt"
	"gophkeeper/internal/client/service"
	"gophkeeper/internal/model"
	"strconv"

	"github.com/spf13/cobra"
)

type CardHandler struct {
	CardService       service.CardService
	UserMasterService service.UserMasterService
}

func NewCardHandler(cardService service.CardService, u service.UserMasterService) *CardHandler {
	return &CardHandler{
		CardService:       cardService,
		UserMasterService: u,
	}

}
func (ch *CardHandler) GetCommands() []*cobra.Command {
	addCardCmd := &cobra.Command{
		Use:     "card-add",
		Short:   "Добавить новую банковскую карту",
		Example: `  todo card-add --title "Название карты" --number "1234 1234 1234 1234" --expiry "12/23" --card_holder "Имя и фамилия" --cvv "123"`,
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
	addCardCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный)")
	addCardCmd.MarkFlagRequired("masterPass")

	getCardCmd := &cobra.Command{
		Use:     "card-get",
		Short:   "Получить список банковских карт",
		Example: `  todo card-get --masterPass "password"`,
		Run:     ch.getCardRun,
	}
	getCardCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный)")
	getCardCmd.MarkFlagRequired("masterPass")

	deleteCardCmd := &cobra.Command{
		Use:     "card-delete",
		Short:   "Удалить банковскую карту",
		Example: `  todo card-delete --masterPass "password" `,
		Run:     ch.deleteCardRun,
	}
	deleteCardCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный)")
	deleteCardCmd.MarkFlagRequired("masterPass")
	deleteCardCmd.Flags().StringP("id", "i", "", "ID карты (Обязательный)")
	deleteCardCmd.MarkFlagRequired("id")

	updateCardCmd := &cobra.Command{
		Use:     "card-update",
		Short:   "Добавить новую банковскую карту",
		Example: `  todo card-update --title "Название карты" --number "1234 1234 1234 1234" --expiry "12/23" --card_holder "Имя и фамилия" --cvv "123"`,
		Run:     ch.updateCardRun,
	}
	updateCardCmd.Flags().StringP("id", "i", "", "ID карты (обязательное)")
	updateCardCmd.MarkFlagRequired("id")
	updateCardCmd.Flags().StringP("title", "t", "", "Произвольное название карты (обязательное)")
	updateCardCmd.MarkFlagRequired("title")
	updateCardCmd.Flags().StringP("number", "n", "", "Номер карты (обязательное)")
	updateCardCmd.MarkFlagRequired("number")
	updateCardCmd.Flags().StringP("expiry", "e", "", "Срок действия карты (обязательное)")
	updateCardCmd.MarkFlagRequired("expiry")
	updateCardCmd.Flags().StringP("cvv", "c", "", "cvv (обязательное)")
	updateCardCmd.MarkFlagRequired("cvv")
	updateCardCmd.Flags().StringP("card_holder", "o", "", "Имя держателя карты (обязательное)")
	updateCardCmd.MarkFlagRequired("card_holder")
	updateCardCmd.Flags().StringP("masterPass", "m", "", "Пароль для шифрования (Обязательный)")
	updateCardCmd.MarkFlagRequired("masterPass")

	return []*cobra.Command{
		addCardCmd,
		getCardCmd,
		deleteCardCmd,
		updateCardCmd,
	}
}

func (ch *CardHandler) deleteCardRun(cmd *cobra.Command, _ []string) {
	masterPass, err := cmd.Flags().GetString("masterPass")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}
	ch.CardService.SetMasterPass(masterPass)
	err = ch.UserMasterService.Master(masterPass)
	if err != nil {
		fmt.Println("❌ Ошибка проверки мастер пароля:", err)
		return
	}

	id, err := cmd.Flags().GetString("id")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите порядковый номер карты (--id или -i)")
		return
	}

	err = ch.CardService.DeleteCard(id)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	fmt.Printf("✅ Карта успешно удалена\n")
}

func (ch *CardHandler) getCardRun(cmd *cobra.Command, _ []string) {
	masterPass, err := cmd.Flags().GetString("masterPass")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}
	ch.CardService.SetMasterPass(masterPass)

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

	for _, card := range cards {
		fmt.Printf("📋 Карта # %d\n", card.ID)
		fmt.Printf("Title: %s\n", card.Title)
		fmt.Printf("Number: %s\n", card.Number)
		fmt.Printf("Expiry: %s\n", card.Expiry)
		fmt.Printf("Card Holder: %s\n", card.CardHolder)
		fmt.Printf("CVV: %s\n", card.CVV)
		fmt.Println()
	}
}

func (ch *CardHandler) addCardRun(cmd *cobra.Command, _ []string) {
	masterPass, err := cmd.Flags().GetString("masterPass")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}
	ch.CardService.SetMasterPass(masterPass)

	err = ch.UserMasterService.Master(masterPass)
	if err != nil {
		fmt.Println("❌ Ошибка проверки мастер пароля:", err)
		return
	}

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

func (ch *CardHandler) updateCardRun(cmd *cobra.Command, _ []string) {
	masterPass, err := cmd.Flags().GetString("masterPass")
	if err != nil || masterPass == "" {
		fmt.Println("❌ Ошибка: укажите пароль (--masterPass или -m)")
		return
	}
	ch.CardService.SetMasterPass(masterPass)

	err = ch.UserMasterService.Master(masterPass)
	if err != nil {
		fmt.Println("❌ Ошибка проверки мастер пароля:", err)
		return
	}

	card := model.Card{}
	str, err := cmd.Flags().GetString("title")
	if err != nil || str == "" {
		fmt.Println("❌ Ошибка: укажите произвольное название карты (--title или -t)")
		return
	}
	card.Title = str

	id, err := cmd.Flags().GetString("id")
	if err != nil || id == "" {
		fmt.Println("❌ Ошибка: укажите произвольное название карты (--id или -i)")
		return
	}
	card.ID, err = strconv.Atoi(id)
	if err != nil {
		fmt.Println("❌ Ошибка: не корректный формат ID карты")
		return
	}

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

	err = ch.CardService.UpdateCard(card)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	fmt.Printf("✅ Запись успешно обновлена\n")

}
