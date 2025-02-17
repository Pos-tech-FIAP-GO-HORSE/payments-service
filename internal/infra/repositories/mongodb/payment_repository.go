package mongodb

import (
	"context"
	"github.com/Pos-tech-FIAP-GO-HORSE/payments-service/internal/core/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection) *PaymentRepository {
	return &PaymentRepository{collection: collection}
}

func (u *PaymentRepository) Save(ctx context.Context, payment *entities.Payment) error {
	_, err := u.collection.InsertOne(ctx, payment)
	if err != nil {
		return err
	}
	return nil
}

func (u *PaymentRepository) FindByID(ctx context.Context, id string) (*entities.Payment, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := u.collection.FindOne(ctx, bson.M{"_id": objectID})
	if err := result.Err(); err != nil {
		return nil, err
	}

	var payment entities.Payment
	if err := result.Decode(&payment); err != nil {
		return nil, err
	}

	return &payment, nil
}
