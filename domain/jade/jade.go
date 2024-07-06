package jade_domain

import (
	"clean-architecture/internal/blockchain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

const (
	CollectionJade = "jade"
)

// Jade is a currency object, and this collection just only create and update from second interaction
type Jade struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount *blockchain.Block  `bson:"amount" json:"amount"` // amount data is a blockchain
}

type JadeBlockchain struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount int64              `bson:"amount" json:"amount"` // amount data is a blockchain
}

type JadeTransaction struct {
	Jade            Jade   `bson:"jade" json:"jade"`
	TransactionType string `bson:"transaction_type" json:"transaction_type"`
}

type IJadeRepository interface {
	FetchJadeInUser(ctx context.Context, userID primitive.ObjectID) (JadeBlockchain, error)
	CreateCurrencyInUser(ctx context.Context, userID primitive.ObjectID, data int64) error
	UpdateCurrencyInUser(ctx context.Context, userID primitive.ObjectID, data int64) error
	Rank(ctx context.Context) ([]JadeBlockchain, error)
}
