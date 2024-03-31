package mean_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionMean = "mean"
)

type Mean struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`
	Description  string             `bson:"description" json:"description"`
	Example      string             `bson:"example" json:"example"`
	VietSub      string             `bson:"viet_sub" json:"viet_sub"`
	SynonymID    string             `bson:"synonym" json:"synonym"`
	AntonymID    string             `bson:"antonyms" json:"antonyms"`
}

type Response struct {
	Mean []Mean `bson:"data" json:"data"`
}

type IMeanRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	CreateOne(ctx context.Context, mean *Mean, fieldOfIT string) error
	UpdateOne(ctx context.Context, meanID string, mean Mean) error
	UpsertOne(ctx context.Context, id string, mean *Mean) (Mean, error)
	DeleteOne(ctx context.Context, meanID string) error
}
