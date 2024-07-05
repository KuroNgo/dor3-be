package markList_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	NameList    string             `bson:"name_list" json:"name_list"`
	Description string             `bson:"description" json:"description"`
}

type IMarkListUseCase interface {
	FetchManyByUser(ctx context.Context, user primitive.ObjectID) (Response, error)
	FetchByIdByUser(ctx context.Context, user primitive.ObjectID, id string) (MarkList, error)

	CreateOneByUser(ctx context.Context, user primitive.ObjectID, markList *MarkList) error
	UpdateOneByUser(ctx context.Context, user primitive.ObjectID, markList *MarkList) (*mongo.UpdateResult, error)
	DeleteOneByUser(ctx context.Context, user primitive.ObjectID, markListID string) error

	FetchManyByAdmin(ctx context.Context) (Response, error)
	FetchByIdByAdmin(ctx context.Context, id string) (MarkList, error)
}
