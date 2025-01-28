package usecases

import (
	"context"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/entities"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/interfaces"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
)

type GeneratePayment struct {
	GeneratorPayment interfaces.IGeneratorPayment
}

func NewGeneratePayment(generatorPayment interfaces.IGeneratorPayment) *GeneratePayment {
	return &GeneratePayment{
		GeneratorPayment: generatorPayment,
	}
}

func (uc *GeneratePayment) Execute(ctx context.Context, input entities.Input) (*dto.ResponseCreatePayment, error) {
	paymentInfos, err := uc.GeneratorPayment.GeneratePaymentToOrder(ctx, input.Amount, input.OrderID)
	if err != nil {
		return nil, err
	}
	return paymentInfos, nil
}
