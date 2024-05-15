package mark_vacabulary_usecase

import (
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	"context"
	"time"
)

type markVocabularyUseCase struct {
	markVocabularyRepository mark_vocabulary_domain.IMarkToFavouriteRepository
	contextTimeout           time.Duration
}

func NewMarkVocabularyUseCase(markVocabularyRepository mark_vocabulary_domain.IMarkToFavouriteRepository, timeout time.Duration) mark_vocabulary_domain.IMarkToFavouriteRepository {
	return &markVocabularyUseCase{
		markVocabularyRepository: markVocabularyRepository,
		contextTimeout:           timeout,
	}
}

func (m *markVocabularyUseCase) FetchManyByMarkList(ctx context.Context, markListId string) (mark_vocabulary_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	markVocabulary, err := m.markVocabularyRepository.FetchManyByMarkList(ctx, markListId)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}

	return markVocabulary, err
}

func (m *markVocabularyUseCase) FetchManyByMarkListIDAndUserId(ctx context.Context, markListId string, userId string) (mark_vocabulary_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	markVocabulary, err := m.markVocabularyRepository.FetchManyByMarkListIDAndUserId(ctx, markListId, userId)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}

	return markVocabulary, err
}

func (m *markVocabularyUseCase) FetchManyByMarkListID(ctx context.Context, markListId string) ([]mark_vocabulary_domain.MarkToFavourite, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	markVocabulary, err := m.markVocabularyRepository.FetchManyByMarkListID(ctx, markListId)
	if err != nil {
		return nil, err
	}

	return markVocabulary, err
}

func (m *markVocabularyUseCase) UpdateOne(ctx context.Context, markVocabularyID string, markVocabulary mark_vocabulary_domain.MarkToFavourite) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.markVocabularyRepository.UpdateOne(ctx, markVocabularyID, markVocabulary)
	if err != nil {
		return err
	}

	return nil
}

func (m *markVocabularyUseCase) CreateOne(ctx context.Context, markVocabulary *mark_vocabulary_domain.MarkToFavourite) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.markVocabularyRepository.CreateOne(ctx, markVocabulary)
	if err != nil {
		return err
	}

	return nil
}

func (m *markVocabularyUseCase) DeleteOne(ctx context.Context, markVocabularyID string) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.markVocabularyRepository.DeleteOne(ctx, markVocabularyID)
	if err != nil {
		return err
	}

	return nil
}
