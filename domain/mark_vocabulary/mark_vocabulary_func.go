package mark_vocabulary_domain

import "context"

type IMarkToFavouriteUseCase interface {
	FetchManyByMarkListIDAndUserId(ctx context.Context, markListId string, userId string) (Response, error)

	CreateOne(ctx context.Context, markVocabulary *MarkToFavourite) error
	UpdateOne(ctx context.Context, markVocabularyID string, markVocabulary MarkToFavourite) error
	DeleteOne(ctx context.Context, markVocabularyID string) error
}
