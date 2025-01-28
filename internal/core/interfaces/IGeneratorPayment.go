package interfaces

import (
	"context"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
)

type IGeneratorPayment interface {
	GeneratePaymentToOrder(ctx context.Context, amount float64) (*dto.ResponseCreatePayment, error)
	GetPaymentStatus(ctx context.Context, paymentId int) (*dto.ResponseStatusPayment, error)
}
