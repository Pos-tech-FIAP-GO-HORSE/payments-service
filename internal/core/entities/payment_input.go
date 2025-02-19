package entities

import (
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
	"time"
)

type Input struct {
	Amount   float64 `json:"amount"`
	OrderID  string  `json:"order_id"`
	PublicID string  `json:"public_id"`
}

type ResponseCreatePayment struct {
	Message dto.ResponseCreatePayment `json:"result"`
	Error   string                    `json:"error,omitempty"`
}

type ResponseStatusPayment struct {
	Message dto.ResponseStatusPayment `json:"result"`
	Error   string                    `json:"error,omitempty"`
}

type Payment struct {
	Amount    float64 `json:"amount"`
	OrderID   string  `json:"order_id"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	PublicID  string  `json:"public_id"`
}

func NewPayment(amount float64, orderID, status string, publicID string) *Payment {
	return &Payment{
		Amount:    amount,
		OrderID:   orderID,
		Status:    status,
		CreatedAt: time.Now().Format(time.RFC3339),
		PublicID:  publicID,
	}
}
