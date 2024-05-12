package markList_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	NameList    string             `bson:"name_list" json:"name_list"`
	Description string             `bson:"description" json:"description"`
}

type IMarkListUseCase interface {
	FetchManyByUserID(ctx context.Context, userId string) (Response, error)
	FetchMany(ctx context.Context) (Response, error)
	FetchById(ctx context.Context, id string) (MarkList, error)

	CreateOne(ctx context.Context, markList *MarkList) error
	UpdateOne(ctx context.Context, markList *MarkList) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, markListID string) error
}
