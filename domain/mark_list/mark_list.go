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
	FetchManyByUserID(ctx context.Context, userId string) (Response, error)
	FetchMany(ctx context.Context) (Response, error)
	FetchById(ctx context.Context, id string) (MarkList, error)

	CreateOne(ctx context.Context, markList *MarkList) error
	UpdateOne(ctx context.Context, markList *MarkList) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, markListID string) error

	Statistics(ctx context.Context) (Statistics, error)
}
