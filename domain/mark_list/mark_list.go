package markList_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionMarkList = "mark_list"
)

type MarkList struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	NameList    string             `bson:"name_list" json:"name_list"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	WhoCreated  string             `bson:"who_created" json:"who_created"`
}

type Response struct {
	Statistics Statistics `json:"statistics" bson:"statistics"`
	MarkList   []MarkList `json:"mark_list" bson:"mark_list"`
}

type Statistics struct {
	Total           int64 `json:"total"`
	CountVocabulary int64 `json:"count_vocabulary"`
}

type IMarkListRepository interface {
	FetchManyByUser(ctx context.Context, user primitive.ObjectID) (Response, error)
	FetchByIdByUser(ctx context.Context, user primitive.ObjectID, id string) (MarkList, error)
	StatisticsIndividual(ctx context.Context, user primitive.ObjectID) (Statistics, error)

	CreateOneByUser(ctx context.Context, user primitive.ObjectID, markList *MarkList) error
	UpdateOneByUser(ctx context.Context, user primitive.ObjectID, markList *MarkList) (*mongo.UpdateResult, error)
	DeleteOneByUser(ctx context.Context, user primitive.ObjectID, markListID string) error

	FetchManyByAdmin(ctx context.Context) (Response, error)
	Statistics(ctx context.Context) (Statistics, error)
}
