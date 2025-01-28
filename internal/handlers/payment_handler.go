package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/entities"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/interfaces"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/usecases"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PaymentHandler struct {
	GeneratePaymentUseCase *usecases.GeneratePayment
	StatusPaymentUseCase   *usecases.StatusPayment
}

func NewPaymentHandler(generatorPayment interfaces.IGeneratorPayment) *PaymentHandler {
	return &PaymentHandler{
		GeneratePaymentUseCase: usecases.NewGeneratePayment(generatorPayment),
		StatusPaymentUseCase:   usecases.NewStatusPayment(generatorPayment),
	}
}

func (h *PaymentHandler) HandleCreatePayment(ctx context.Context, snsEvent events.SNSEvent) error {
	for _, record := range snsEvent.Records {
		sns := record.SNS
		var input entities.Input

		err := json.Unmarshal([]byte(sns.Message), &input)
		if err != nil {
			log.Fatalf("Erro ao fazer unmarshal: %v", err)
		}
		_, err = h.GeneratePaymentUseCase.Execute(ctx, input)
		if err != nil {
			log.Fatalf("Erro ao gerar pagamento: %v", err)
		}
		log.Println("Pagamento gerado com sucesso")
	}
	return nil
}

//// CreatePayment godoc
//// @Summary      Create a new payment
//// @Description  Add a new payment to order
//// @Tags         Payments
//// @Accept       json
//// @Produce      json
//// @Param        create_payment   body      payment.Input  true  "Payment Data"
//// @Success      200     {object}  ResponseCreatePayment
//// @Failure      400     {object}  ResponseCreatePayment
//// @Failure      500     {object}  ResponseCreatePayment
//// @Router       /api/v1/payments [post]
//func (h *PaymentHandler) HandleCreatePayment(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	var input entities.Input
//	// Bind do corpo da requisição para a estrutura de dados de entrada
//	err := json.Unmarshal([]byte(request.Body), &input)
//	if err != nil {
//		return events.APIGatewayProxyResponse{
//			StatusCode: http.StatusBadRequest,
//			Body:       fmt.Sprintf("Erro ao deserializar o corpo: %v", err),
//		}, nil
//	}
//
//	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
//	defer cancel()
//
//	paymentInfos, err := h.GeneratePaymentUseCase.Execute(ctx, input)
//	if err != nil {
//		return events.APIGatewayProxyResponse{
//			StatusCode: http.StatusInternalServerError,
//			Body:       fmt.Sprintf("Erro ao gerar pagamento: %v", err),
//		}, nil
//	}
//
//	// Resposta de sucesso
//	return events.APIGatewayProxyResponse{
//		StatusCode: http.StatusOK,
//		Body: fmt.Sprintf(`{"message":"payment created successfully", "paymentQRCode":"%s", "paymentId":"%s"}`,
//			paymentInfos.QRCode, paymentInfos.ID),
//	}, nil
//}

// GetStatusPayment StatusPayment godoc
// @Summary      Get a payment status
// @Description  Get a payment order status
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id   path     string  true  "Payment ID"
// @Success      200     {object}  ResponseStatusPayment
// @Failure      400     {object}  ResponseStatusPayment
// @Failure      500     {object}  ResponseStatusPayment
// @Router       /api/v1/payments/{id} [get]
func (h *PaymentHandler) HandleGetStatusPayment(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	param := request.PathParameters["id"]
	if param == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Payment ID is required"}`,
		}, nil
	}

	paymentID, err := strconv.Atoi(param)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid Payment ID"}`,
		}, nil
	}

	// Lógica de status de pagamento com timeout
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	paymentStatusResponse, err := h.StatusPaymentUseCase.Execute(ctx, paymentID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Erro ao obter status do pagamento: %v", err),
		}, nil
	}

	// Retorno com status do pagamento
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"result":"%s"}`, paymentStatusResponse),
	}, nil
}
