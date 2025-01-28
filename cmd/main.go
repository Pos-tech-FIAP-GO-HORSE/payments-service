package main

import (
	"context"
	"fmt"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/handlers"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/client"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	mercadopagoclient "github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"log"
)

func main() {

	var (
		tokenMP = "TEST-2373946154784631-101516-50ff7f4dcdff3aec43372568c77990e3-175794680"
		//os.Getenv("TOKEN_MERCADO_PAGO")
	)

	// Client
	cfg, err := mercadopagoclient.New(tokenMP)
	if err != nil {
		log.Fatalf("Erro ao criar configuração: %v", err)
	}

	mpClient := payment.NewClient(cfg)
	paymentClient := client.NewGeneratorPayment(mpClient)

	// Handler
	paymentHandler := handlers.NewPaymentHandler(paymentClient)

	lambda.Start(func(ctx context.Context, event interface{}) (interface{}, error) {
		switch e := event.(type) {
		case events.APIGatewayProxyRequest:
			return paymentHandler.HandleGetStatusPayment(ctx, e)
		case events.SNSEvent:
			return nil, paymentHandler.HandleCreatePayment(ctx, e)
		default:
			return nil, fmt.Errorf("unknown event type: %T", e)
		}
	})
}
