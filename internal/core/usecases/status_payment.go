package usecases

import (
	"context"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/interfaces"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
)

type StatusPayment struct {
	GeneratorPayment interfaces.IGeneratorPayment
}

func NewStatusPayment(generatorPayment interfaces.IGeneratorPayment) *StatusPayment {
	return &StatusPayment{
		GeneratorPayment: generatorPayment,
	}
}

func (uc *StatusPayment) Execute(ctx context.Context, paymentId int) (*dto.ResponseStatusPayment, error) {
	paymentInfos, err := uc.GeneratorPayment.GetPaymentStatus(ctx, paymentId)
	if err != nil {
		return nil, err
	}
	return paymentInfos, nil
}
