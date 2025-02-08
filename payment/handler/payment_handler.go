package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kaweel/workshop-tdd/payment/service"
)

type PaymentHandler interface {
	Payment() http.HandlerFunc
}

type paymentHandler struct {
	p service.Service
}

func NewHandler(p service.Service) PaymentHandler {
	return &paymentHandler{
		p: p,
	}
}

func (h *paymentHandler) Payment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.RequestPayment

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = h.p.Payment(req)
		if err != nil {
			var httpstatus int
			switch err.Error() {
			case "invalid payment channel",
				"order status is not request payment",
				"customer status is not active",
				"customer amount is not enough",
				"merchant status is not active":
				httpstatus = http.StatusUnprocessableEntity
			default:
				httpstatus = http.StatusInternalServerError
			}

			http.Error(w, fmt.Sprintf("Failed payment : %v", err.Error()), httpstatus)
			return
		}

	}
}
