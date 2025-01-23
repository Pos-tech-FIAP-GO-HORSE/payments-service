package infra

import (
	"github.com/mercadopago/sdk-go/pkg/payment"
)

type GeneratorPayment struct {
	Client payment.Client
}

func NewPaymentClient(client payment.Client) *Payment {
	return &Payment{
		Client: client,
	}
}
