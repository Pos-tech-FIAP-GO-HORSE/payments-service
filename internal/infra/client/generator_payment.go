package client

import (
	"context"
	"fmt"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

type GeneratorPayment struct {
	Client payment.Client
}

func NewGeneratorPayment(client payment.Client) *GeneratorPayment {
	return &GeneratorPayment{
		Client: client,
	}
}

func (p *GeneratorPayment) GeneratePaymentToOrder(ctx context.Context, amount float64) (*dto.ResponseCreatePayment, error) {
	request := payment.Request{
		TransactionAmount: amount,
		PaymentMethodID:   "pix",
	}

	response, err := p.Client.Create(ctx, request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	responsePayment := &dto.ResponseCreatePayment{
		QRCode: response.PointOfInteraction.TransactionData.QRCode,
		ID:     response.ID,
	}

	return responsePayment, nil
}

func (p *GeneratorPayment) GetPaymentStatus(ctx context.Context, paymentId int) (*dto.ResponseStatusPayment, error) {
	response, err := p.Client.Get(ctx, paymentId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	responseStatusPayment := &dto.ResponseStatusPayment{
		ID:                response.ID,
		Status:            response.Status,
		StatusDetail:      response.StatusDetail,
		TransactionAmount: response.TransactionAmount,
		PaymentMethodId:   response.PaymentMethodID,
	}

	return responseStatusPayment, nil
}
