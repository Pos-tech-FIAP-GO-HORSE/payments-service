package interfaces

import (
	"context"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/entities"
)

type IPaymentRepository interface {
	Save(ctx context.Context, payment *entities.Payment) error
	FindByID(ctx context.Context, id string) (*entities.Payment, error)
}
