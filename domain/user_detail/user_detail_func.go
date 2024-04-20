package user_detail_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Specialize string             `bson:"specialize"  json:"specialize"`
}

type IUserDetailUseCase interface {
	FetchByUserID(c context.Context) (UserDetail, error)
	Create(c context.Context, user UserDetail) error
	Update(ctx context.Context, user *UserDetail) (*mongo.UpdateResult, error)
	DeleteOne(c context.Context, userID string) error
}
