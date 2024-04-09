package markList_repository

import "context"

type Input struct {
	NameList    string `bson:"name_list" json:"name_list"`
	Description string `bson:"description" json:"description"`
}

type IMarkListUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	UpdateOne(ctx context.Context, markListID string, markList MarkList) error
	CreateOne(ctx context.Context, markList *MarkList) error
	UpsertOne(c context.Context, id string, markList *MarkList) (Response, error)
	DeleteOne(ctx context.Context, markListID string) error
}
