package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/entities"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/usecases"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
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

func NewPaymentHandler(generatePayment *usecases.GeneratePayment, statusPayment *usecases.StatusPayment) *PaymentHandler {
	return &PaymentHandler{
		GeneratePaymentUseCase: generatePayment,
		StatusPaymentUseCase:   statusPayment,
	}
}

func (h *PaymentHandler) HandleCreatePayment(ctx context.Context, snsEvent events.SNSEvent) (*dto.ResponseCreatePayment, error) {
	for _, record := range snsEvent.Records {
		sns := record.SNS

		var record string
		retorno, _ := json.Marshal(record)
		log.Println("Message received from SNS:", string(retorno))

		var rawMessage string
		err := json.Unmarshal([]byte(sns.Message), &rawMessage)
		if err != nil {
			log.Fatal("Failed to deserialize initial JSON:", err)
			return nil, err
		}

		var input entities.Input
		err = json.Unmarshal([]byte(rawMessage), &input)
		if err != nil {
			log.Fatal("Failed to deserialize final JSON:", err)
			return nil, err
		}

		paymentGenerated, err := h.GeneratePaymentUseCase.Execute(ctx, input)
		if err != nil {
			log.Fatalf("Error while generating payment: %v", err)
			return nil, err
		}
		log.Println("Payment created successfully:", paymentGenerated)
		return paymentGenerated, nil
	}
	return nil, nil
}

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
