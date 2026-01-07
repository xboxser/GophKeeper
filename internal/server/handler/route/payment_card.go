package route

import (
	"encoding/json"
	"gophkeeper/internal/server/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PaymentCardHandler struct {
	PaymentCardService service.PaymentCardService
}

func NewPaymentCardHandler(paymentCardService service.PaymentCardService) *PaymentCardHandler {
	return &PaymentCardHandler{
		PaymentCardService: paymentCardService,
	}
}

func (ph *PaymentCardHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", ph.getPaymentCards)
	return r
}
func (ph *PaymentCardHandler) Pattern() string {
	return "/api/payment_cards"
}

func (ph *PaymentCardHandler) getPaymentCards(w http.ResponseWriter, r *http.Request) {
	paymentCards, err := ph.PaymentCardService.GetPaymentCards()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paymentCards)
}
