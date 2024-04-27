package mark_vocabulary_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionMark = "add_to_favourite"
)

type MarkToFavourite struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId       primitive.ObjectID `bson:"user_id" json:"user_id"`
	MarkListID   primitive.ObjectID `bson:"mark_list_id" json:"mark_list_id"`
	VocabularyID primitive.ObjectID `bson:"mark_vocabulary_id" json:"mark_vocabulary_id"`
}

type Response struct {
	MarkToFavourite []MarkToFavourite
}

type IMarkToFavouriteRepository interface {
	FetchManyByMarkListIDAndUserId(ctx context.Context, markListId string, userId string) (Response, error)

	CreateOne(ctx context.Context, markVocabulary *MarkToFavourite) error
	UpdateOne(ctx context.Context, markVocabularyID string, markVocabulary MarkToFavourite) error
	DeleteOne(ctx context.Context, markVocabularyID string) error
}
