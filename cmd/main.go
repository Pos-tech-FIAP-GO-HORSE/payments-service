package main

import (
	"context"
	"encoding/json"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/usecases"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/handlers"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/client"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/publisher"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/repositories/mongodb"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	mercadopagoclient "github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {

	var (
		dbName  = "pos-tech-fiap"
		dbURI   = "mongodb+srv://admin:admin123@payment-cluster.sl4mh.mongodb.net/pos-tech-fiap?retryWrites=true&w=majority"
		tokenMP = "TEST-2373946154784631-101516-50ff7f4dcdff3aec43372568c77990e3-175794680"
		//os.Getenv("TOKEN_MERCADO_PAGO")
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Client
	cfg, err := mercadopagoclient.New(tokenMP)
	if err != nil {
		log.Fatalf("Erro ao criar configuração: %v", err)
	}
	mpClient := payment.NewClient(cfg)
	paymentClient := client.NewGeneratorPayment(mpClient)

	// Repository
	mongoClient, err := mongodb.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalf("error to connect to database: %v", err)
	}

	database := mongoClient.Database(dbName)
	paymentCollection := database.Collection("payments")
	log.Println(paymentCollection.Name())

	paymentRepository := mongodb.NewPaymentRepository(paymentCollection)

	// Publisher
	messagePublisher, err := publisher.NewSNSService("arn:aws:sns:us-east-1:537124948968:payment-created-service-sns-topic")
	if err != nil {
		log.Fatalf("error to create SNS service: %v", err)
	}

	// UseCase
	generatePaymentUseCase := usecases.NewGeneratePayment(paymentClient, paymentRepository, messagePublisher)
	statusPaymentUseCase := usecases.NewStatusPayment(paymentClient)

	// Handler
	paymentHandler := handlers.NewPaymentHandler(generatePaymentUseCase, statusPaymentUseCase)

	lambda.Start(func(ctx context.Context, event json.RawMessage) (interface{}, error) {
		var apiGatewayEvent events.APIGatewayProxyRequest
		var snsEvent events.SNSEvent

		if err := json.Unmarshal(event, &apiGatewayEvent); err == nil && apiGatewayEvent.Resource != "" {
			return paymentHandler.HandleGetStatusPayment(ctx, apiGatewayEvent)
		}

		if err := json.Unmarshal(event, &snsEvent); err == nil && len(snsEvent.Records) > 0 {
			return paymentHandler.HandleCreatePayment(ctx, snsEvent)
		}

		log.Fatal("Event not supported")
		return nil, nil
	})
}
