package mark_vocabulary_domain

import (
	mark_list_domain "clean-architecture/domain/mark_list"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionMark = "add_to_favourite"
)

type MarkToFavourite struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserId       primitive.ObjectID `bson:"user_id" json:"user_id"`
	MarkListID   primitive.ObjectID `bson:"mark_list_id" json:"mark_list_id"`
	VocabularyID primitive.ObjectID `bson:"mark_vocabulary_id" json:"mark_vocabulary_id"`
}

type MarkToFavouriteResponse struct {
	ID         primitive.ObjectID           `bson:"_id,omitempty" json:"_id"`
	UserId     primitive.ObjectID           `bson:"user_id" json:"user_id"`
	MarkList   mark_list_domain.MarkList    `bson:"mark_list" json:"mark_list"`
	Vocabulary vocabulary_domain.Vocabulary `bson:"vocabulary" json:"vocabulary"`
}

type Response struct {
	Total                   int                       `json:"total" bson:"total"`
	MarkToFavouriteResponse []MarkToFavouriteResponse `json:"mark_to_favourite" bson:"mark_to_favourite"`
}

type IMarkToFavouriteRepository interface {
	FetchManyByMarkListIDAndUserId(ctx context.Context, markListId string, userId string) (Response, error)
	FetchManyByMarkList(ctx context.Context, markListId string) (Response, error)
	FetchManyByMarkListID(ctx context.Context, markListId string) ([]MarkToFavourite, error)

	CreateOne(ctx context.Context, markVocabulary *MarkToFavourite) error
	UpdateOne(ctx context.Context, markVocabularyID string, markVocabulary MarkToFavourite) error
	DeleteOne(ctx context.Context, markVocabularyID string) error
}
