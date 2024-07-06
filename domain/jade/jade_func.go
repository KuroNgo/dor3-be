package jade_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IJadeUseCase interface {
	FetchJadeInUser(ctx context.Context, userID primitive.ObjectID) (JadeBlockchain, error)
	CreateCurrencyInUser(ctx context.Context, userID primitive.ObjectID, data int64) error
	UpdateCurrencyInUser(ctx context.Context, userID primitive.ObjectID, data int64) error
	Rank(ctx context.Context) ([]JadeBlockchain, error)
}
