package markList_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	NameList    string `bson:"name_list" json:"name_list"`
	Description string `bson:"description" json:"description"`
}

type IMarkListUseCase interface {
	FetchManyByUserID(ctx context.Context, userId string) (Response, error)

	CreateOne(ctx context.Context, markList *MarkList) error
	UpdateOne(ctx context.Context, markList *MarkList) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, markListID string) error
}
