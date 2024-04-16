package markList_domain

import "context"

type Input struct {
	NameList    string `bson:"name_list" json:"name_list"`
	Description string `bson:"description" json:"description"`
}

type IMarkListUseCase interface {
	FetchManyByUserID(ctx context.Context, userId string) (Response, error)

	CreateOne(ctx context.Context, markList *MarkList) error
	UpdateOne(ctx context.Context, markListID string, markList MarkList) error
	DeleteOne(ctx context.Context, markListID string) error
}
