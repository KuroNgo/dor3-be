package jade_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionJade = "jade"
)

// Jade is a currency object, and this collection just only create and update from second interaction
type Jade struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount int64              `bson:"amount" json:"amount"` // amount data is a blockchain
}

type JadeTransaction struct {
	Jade            Jade   `bson:"jade" json:"jade"`
	TransactionType string `bson:"transaction_type" json:"transaction_type"`
}
