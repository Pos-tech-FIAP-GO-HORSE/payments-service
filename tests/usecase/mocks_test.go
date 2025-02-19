package usecases

import (
	"context"

	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
	"github.com/stretchr/testify/mock"
)

type MockGeneratorPayment struct {
	mock.Mock
}

func (m *MockGeneratorPayment) GetPaymentStatus(ctx context.Context, paymentId int) (*dto.ResponseStatusPayment, error) {
	args := m.Called(ctx, paymentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ResponseStatusPayment), args.Error(1)
}

func (m *MockGeneratorPayment) GeneratePaymentToOrder(ctx context.Context, amount float64) (*dto.ResponseCreatePayment, error) {
	args := m.Called(ctx, amount)
	return args.Get(0).(*dto.ResponseCreatePayment), args.Error(1)
}
