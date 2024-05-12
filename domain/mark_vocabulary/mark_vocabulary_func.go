package mark_vocabulary_domain

import (
	"context"
)

type Input struct {
	MarkListID   string `bson:"mark_list_id" json:"mark_list_id"`
	VocabularyID string `bson:"mark_vocabulary_id" json:"mark_vocabulary_id"`
}

type IMarkToFavouriteUseCase interface {
	FetchManyByMarkListIDAndUserId(ctx context.Context, markListId string, userId string) (Response, error)
	FetchManyByMarkList(ctx context.Context, markListId string) (Response, error)

	CreateOne(ctx context.Context, markVocabulary *MarkToFavourite) error
	UpdateOne(ctx context.Context, markVocabularyID string, markVocabulary MarkToFavourite) error
	DeleteOne(ctx context.Context, markVocabularyID string) error
}
