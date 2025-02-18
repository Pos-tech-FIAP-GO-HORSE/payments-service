package usecases

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/entities"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/usecases"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) FindByID(ctx context.Context, id string) (*entities.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Save(ctx context.Context, payment *entities.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

// Mock de IMessagePublisher
type MockMessagePublisher struct {
	mock.Mock
}

func (m *MockMessagePublisher) Send(ctx context.Context, message dto.MessageData) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func TestExecute_Success(t *testing.T) {
	ctx := context.Background()
	mockGeneratorPayment := new(MockGeneratorPayment)
	mockPaymentRepository := new(MockPaymentRepository)
	mockMessagePublisher := new(MockMessagePublisher)

	useCase := usecases.NewGeneratePayment(mockGeneratorPayment, mockPaymentRepository, mockMessagePublisher)

	input := entities.Input{
		Amount:   100.0,
		OrderID:  "order-123",
		PublicID: uuid.New(),
	}

	mockResponse := &dto.ResponseCreatePayment{
		QRCode: "test-qrcode",
		ID:     12345,
	}

	mockGeneratorPayment.On("GeneratePaymentToOrder", ctx, input.Amount).Return(mockResponse, nil)

	mockMessagePublisher.On("Send", ctx, mock.Anything).Return(nil)

	mockPaymentRepository.On("Save", ctx, mock.Anything).Return(nil)

	response, err := useCase.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "test-qrcode", response.QRCode)
	assert.Equal(t, 12345, response.ID)

	mockGeneratorPayment.AssertExpectations(t)
	mockMessagePublisher.AssertExpectations(t)
	mockPaymentRepository.AssertExpectations(t)
}

func TestExecute_GeneratePaymentError(t *testing.T) {
	ctx := context.Background()
	mockGeneratorPayment := new(MockGeneratorPayment)
	mockPaymentRepository := new(MockPaymentRepository)
	mockMessagePublisher := new(MockMessagePublisher)

	useCase := usecases.NewGeneratePayment(mockGeneratorPayment, mockPaymentRepository, mockMessagePublisher)

	input := entities.Input{
		Amount:   100.0,
		OrderID:  "order-123",
		PublicID: uuid.New(),
	}

	mockGeneratorPayment.On("GeneratePaymentToOrder", ctx, input.Amount).Return((*dto.ResponseCreatePayment)(nil), errors.New("erro ao gerar pagamento"))

	response, err := useCase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "erro ao gerar pagamento", err.Error())

	mockGeneratorPayment.AssertExpectations(t)
	mockMessagePublisher.AssertNotCalled(t, "Send")
	mockPaymentRepository.AssertNotCalled(t, "Save")
}

func TestExecute_MessagePublishError(t *testing.T) {
	ctx := context.Background()
	mockGeneratorPayment := new(MockGeneratorPayment)
	mockPaymentRepository := new(MockPaymentRepository)
	mockMessagePublisher := new(MockMessagePublisher)

	useCase := usecases.NewGeneratePayment(mockGeneratorPayment, mockPaymentRepository, mockMessagePublisher)

	input := entities.Input{
		Amount:   100.0,
		OrderID:  "order-123",
		PublicID: uuid.New(),
	}

	mockResponse := &dto.ResponseCreatePayment{
		QRCode: "test-qrcode",
		ID:     12345,
	}

	mockGeneratorPayment.On("GeneratePaymentToOrder", ctx, input.Amount).Return(mockResponse, nil)

	mockMessagePublisher.On("Send", ctx, mock.Anything).Return(errors.New("erro ao publicar mensagem"))

	response, err := useCase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "erro ao publicar mensagem", err.Error())

	mockGeneratorPayment.AssertExpectations(t)
	mockMessagePublisher.AssertExpectations(t)
	mockPaymentRepository.AssertNotCalled(t, "Save")
}

func TestExecute_SavePaymentError(t *testing.T) {
	ctx := context.Background()
	mockGeneratorPayment := new(MockGeneratorPayment)
	mockPaymentRepository := new(MockPaymentRepository)
	mockMessagePublisher := new(MockMessagePublisher)

	useCase := usecases.NewGeneratePayment(mockGeneratorPayment, mockPaymentRepository, mockMessagePublisher)

	input := entities.Input{
		Amount:   100.0,
		OrderID:  "order-123",
		PublicID: uuid.New(),
	}

	mockResponse := &dto.ResponseCreatePayment{
		QRCode: "test-qrcode",
		ID:     12345,
	}

	mockGeneratorPayment.On("GeneratePaymentToOrder", ctx, input.Amount).Return(mockResponse, nil)

	mockMessagePublisher.On("Send", ctx, mock.Anything).Return(nil)

	mockPaymentRepository.On("Save", ctx, mock.Anything).Return(errors.New("erro ao salvar pagamento"))

	response, err := useCase.Execute(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "erro ao salvar pagamento", err.Error())

	mockGeneratorPayment.AssertExpectations(t)
	mockMessagePublisher.AssertExpectations(t)
	mockPaymentRepository.AssertExpectations(t)
}
