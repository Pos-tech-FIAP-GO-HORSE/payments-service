package entities

import "github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"

type Input struct {
	Amount  float64 `json:"amount"`
	OrderID string  `json:"orderID"`
}

type ResponseCreatePayment struct {
	Message dto.ResponseCreatePayment `json:"result"`
	Error   string                    `json:"error,omitempty"`
}

type ResponseStatusPayment struct {
	Message dto.ResponseStatusPayment `json:"result"`
	Error   string                    `json:"error,omitempty"`
}
