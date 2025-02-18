package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"log"
)

type SNSService struct {
	client   *sns.Client
	topicARN string
}

func NewSNSService(topicARN string) (*SNSService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração AWS: %v", err)
	}

	client := sns.NewFromConfig(cfg)

	return &SNSService{
		client:   client,
		topicARN: topicARN,
	}, nil
}

func (s *SNSService) Send(ctx context.Context, messageData dto.MessageData) error {
	messageJSON, err := json.Marshal(messageData)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %v", err)
	}
	input := &sns.PublishInput{
		Message:  aws.String(string(messageJSON)),
		TopicArn: aws.String(s.topicARN),
	}

	_, err = s.client.Publish(ctx, input)
	if err != nil {
		return fmt.Errorf("erro ao publicar no SNS: %v", err)
	}

	log.Printf("QR Code publicado com sucesso no SNS: %s", messageData.QRCode)
	return nil
}
