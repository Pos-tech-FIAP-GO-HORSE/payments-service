package interfaces

import (
	"context"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/infra/dto"
)

type IMessagePublisher interface {
	Send(ctx context.Context, messageData dto.MessageData) error
}
