package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/usecases"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
	"github.com/stretchr/testify/assert"
)

func TestStatusPayment_Execute_Success(t *testing.T) {
	ctx := context.Background()
	mockGeneratorPayment := new(MockGeneratorPayment)
	statusPaymentUseCase := usecases.NewStatusPayment(mockGeneratorPayment)

	paymentID := 12345
	mockResponse := &dto.ResponseStatusPayment{
		ID:                paymentID,
		Status:            "approved",
		StatusDetail:      "payment approved",
		TransactionAmount: 150.0,
		PaymentMethodId:   "pix",
	}

	mockGeneratorPayment.On("GetPaymentStatus", ctx, paymentID).Return(mockResponse, nil)

	response, err := statusPaymentUseCase.Execute(ctx, paymentID)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, paymentID, response.ID)
	assert.Equal(t, "approved", response.Status)
	assert.Equal(t, "payment approved", response.StatusDetail)
	assert.Equal(t, 150.0, response.TransactionAmount)
	assert.Equal(t, "pix", response.PaymentMethodId)

	mockGeneratorPayment.AssertExpectations(t)
}

func TestStatusPayment_Execute_Error(t *testing.T) {
	ctx := context.Background()
	mockGeneratorPayment := new(MockGeneratorPayment)
	statusPaymentUseCase := usecases.NewStatusPayment(mockGeneratorPayment)

	paymentID := 12345

	mockGeneratorPayment.On("GetPaymentStatus", ctx, paymentID).Return(nil, errors.New("erro ao obter status do pagamento"))

	response, err := statusPaymentUseCase.Execute(ctx, paymentID)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "erro ao obter status do pagamento", err.Error())

	mockGeneratorPayment.AssertExpectations(t)
}
