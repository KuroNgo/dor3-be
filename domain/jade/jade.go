package jade_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type JadeTransaction struct {
	ID              primitive.ObjectID `bson:"_id" json:"_id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	TransactionType string             `bson:"transaction_type" json:"transaction_type"`
	Amount          int64              `bson:"amount" json:"amount"`
}
