package user_detail_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CollectionUserDetail = "user_detail"
)

type UserDetail struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Specialize string             `bson:"specialize"  json:"specialize"`
}

type Statistics struct {
	CountSpecialize int16 `bson:"count_specialize" json:"count_specialize"`
	CountUser       int16 `bson:"count_user" json:"count_user"`
}

type IUserDetailRepository interface {
	FetchByUserID(ctx context.Context) (UserDetail, error)
	Create(ctx context.Context, user UserDetail) error
	Update(ctx context.Context, user *UserDetail) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, userID string) error
}
