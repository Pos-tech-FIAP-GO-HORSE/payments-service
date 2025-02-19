package usecases

import (
	"context"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/entities"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/interfaces"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
)

type GeneratePayment struct {
	GeneratorPayment  interfaces.IGeneratorPayment
	PaymentRepository interfaces.IPaymentRepository
	MessagePublisher  interfaces.IMessagePublisher
}

func NewGeneratePayment(generatorPayment interfaces.IGeneratorPayment, paymentRepository interfaces.IPaymentRepository,
	messagePublisher interfaces.IMessagePublisher) *GeneratePayment {
	return &GeneratePayment{
		GeneratorPayment:  generatorPayment,
		PaymentRepository: paymentRepository,
		MessagePublisher:  messagePublisher,
	}
}

func (uc *GeneratePayment) Execute(ctx context.Context, input entities.Input) (*dto.ResponseCreatePayment, error) {
	paymentInfos, err := uc.GeneratorPayment.GeneratePaymentToOrder(ctx, input.Amount)
	if err != nil {
		return nil, err
	}

	messageData := dto.NewMessageData(paymentInfos.QRCode, paymentInfos.ID, input.PublicID)
	err = uc.MessagePublisher.Send(ctx, *messageData)
	if err != nil {
		return nil, err
	}

	payment := entities.NewPayment(input.Amount, input.OrderID, "pendente", input.PublicID)
	err = uc.PaymentRepository.Save(ctx, payment)
	if err != nil {
		return nil, err
	}

	return paymentInfos, nil
}
