package mean_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mean struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`
	Description  string             `bson:"description" json:"description"`
	Example      string             `bson:"example" json:"example"`
	Synonym      string             `bson:"synonym" json:"synonym"`
	Antonym      string             `bson:"antonyms" json:"antonyms"`
}

type Response struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`
	Description  string             `bson:"description" json:"description"`
	Example      string             `bson:"example" json:"example"`
	Synonym      string             `bson:"synonym" json:"synonym"`
	Antonym      string             `bson:"antonyms" json:"antonyms"`
}

type IMeanRepository interface {
	FetchMany(ctx context.Context) ([]Response, error)
	CreateOne(ctx context.Context, mean *Mean, fieldOfIT string) error
	UpdateOne(ctx context.Context, meanID string, mean Mean) error
	UpsertOne(ctx context.Context, id string, mean *Mean) (Mean, error)
	DeleteOne(ctx context.Context, meanID string) error
}
