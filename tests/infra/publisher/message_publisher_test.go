package publisher_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSNSClient struct {
	mock.Mock
}

func (m *MockSNSClient) Publish(ctx context.Context, input *sns.PublishInput, opts ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, input)
	output, _ := args.Get(0).(*sns.PublishOutput)
	return output, args.Error(1)
}

type SNSServiceMock struct {
	client   *MockSNSClient
	topicARN string
}

func (s *SNSServiceMock) Send(ctx context.Context, messageData dto.MessageData) error {
	messageJSON, err := json.Marshal(messageData)
	if err != nil {
		return errors.New("erro ao serializar mensagem")
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(messageJSON)),
		TopicArn: aws.String(s.topicARN),
	}

	_, err = s.client.Publish(ctx, input)
	if err != nil {
		return errors.New("erro ao publicar no SNS")
	}

	return nil
}

func TestSend_Success(t *testing.T) {
	ctx := context.Background()
	mockSNS := new(MockSNSClient)
	topicARN := "arn:aws:sns:us-east-1:123456789012:MyTopic"

	// Criamos o serviço com o mock
	snsService := &SNSServiceMock{
		client:   mockSNS,
		topicARN: topicARN,
	}

	messageData := dto.MessageData{QRCode: "test-qr-code"}
	messageJSON, _ := json.Marshal(messageData)

	mockSNS.On("Publish", ctx, &sns.PublishInput{
		Message:  aws.String(string(messageJSON)),
		TopicArn: aws.String(topicARN),
	}).Return(&sns.PublishOutput{
		MessageId: aws.String("12345"),
	}, nil)

	err := snsService.Send(ctx, messageData)

	assert.NoError(t, err)
	mockSNS.AssertExpectations(t)
}

func TestSend_SNSPublishSuccess(t *testing.T) {
	ctx := context.Background()
	mockSNS := new(MockSNSClient)
	topicARN := "arn:aws:sns:us-east-1:123456789012:MyTopic"

	snsService := &SNSServiceMock{
		client:   mockSNS,
		topicARN: topicARN,
	}

	messageData := dto.MessageData{QRCode: "test-qr-code"}
	messageJSON, _ := json.Marshal(messageData)

	mockSNS.On("Publish", ctx, &sns.PublishInput{
		Message:  aws.String(string(messageJSON)),
		TopicArn: aws.String(topicARN),
	}).Return((*sns.PublishOutput)(nil), nil)

	err := snsService.Send(ctx, messageData)
	assert.NoError(t, err)
	mockSNS.AssertExpectations(t)
}

func TestSend_SNSPublishError(t *testing.T) {
	ctx := context.Background()
	mockSNS := new(MockSNSClient)
	topicARN := "arn:aws:sns:us-east-1:123456789012:MyTopic"

	snsService := &SNSServiceMock{
		client:   mockSNS,
		topicARN: topicARN,
	}

	messageData := dto.MessageData{QRCode: "test-qr-code"}
	messageJSON, _ := json.Marshal(messageData)

	mockSNS.On("Publish", ctx, &sns.PublishInput{
		Message:  aws.String(string(messageJSON)),
		TopicArn: aws.String(topicARN),
	}).Return((*sns.PublishOutput)(nil), errors.New("SNS publish error"))

	err := snsService.Send(ctx, messageData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao publicar no SNS")
	mockSNS.AssertExpectations(t)
}
