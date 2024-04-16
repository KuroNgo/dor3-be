package markList_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionMarkList = "mark_list"
)

type MarkList struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	NameList    string             `bson:"name_list" json:"name_list"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	WhoCreated  string             `bson:"who_created" json:"who_created"`
}

type Response struct {
	MarkList []MarkList
}

type IMarkListRepository interface {
	FetchManyByUserID(ctx context.Context, userId string) (Response, error)

	CreateOne(ctx context.Context, markList *MarkList) error
	UpdateOne(ctx context.Context, markListID string, markList MarkList) error
	DeleteOne(ctx context.Context, markListID string) error
}
