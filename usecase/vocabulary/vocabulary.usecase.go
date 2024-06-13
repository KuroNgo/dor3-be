package vocabulary_usecase

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type vocabularyUseCase struct {
	vocabularyRepository vocabulary_domain.IVocabularyRepository
	contextTimeout       time.Duration
}

func NewVocabularyUseCase(vocabularyRepository vocabulary_domain.IVocabularyRepository, timeout time.Duration) vocabulary_domain.IVocabularyUseCase {
	return &vocabularyUseCase{
		vocabularyRepository: vocabularyRepository,
		contextTimeout:       timeout,
	}
}

func (v *vocabularyUseCase) FindVocabularyIDByVocabularyConfigInAdmin(ctx context.Context, word string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.FindVocabularyIDByVocabularyConfigInAdmin(ctx, word)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) GetLatestVocabularyInAdmin(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.GetLatestVocabularyInAdmin(ctx)
	if err != nil {
		return nil, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) GetVocabularyByIdInAdmin(ctx context.Context, id string) (vocabulary_domain.Vocabulary, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.GetVocabularyByIdInAdmin(ctx, id)
	if err != nil {
		return vocabulary_domain.Vocabulary{}, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) FetchByIdUnitInAdmin(ctx context.Context, idUnit string) ([]vocabulary_domain.Vocabulary, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.FetchByIdUnitInAdmin(ctx, idUnit)
	if err != nil {
		return nil, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) FetchByWordInBoth(ctx context.Context, word string) (vocabulary_domain.SearchingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.FetchByWordInBoth(ctx, word)
	if err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) FetchByLessonInBoth(ctx context.Context, lessonName string) (vocabulary_domain.SearchingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.FetchByLessonInBoth(ctx, lessonName)
	if err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) FetchManyInBoth(ctx context.Context, page string) (vocabulary_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.FetchManyInBoth(ctx, page)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) GetAllVocabularyInAdmin(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.GetAllVocabularyInAdmin(ctx)
	if err != nil {
		return nil, err
	}

	return vocabulary, nil
}

func (v *vocabularyUseCase) FindUnitIDByUnitLevelInAdmin(ctx context.Context, unitLevel int, fieldOfIT string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	unitID, err := v.vocabularyRepository.FindUnitIDByUnitLevelInAdmin(ctx, unitLevel, fieldOfIT)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return unitID, err
}

func (v *vocabularyUseCase) CreateOneByNameUnitInAdmin(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()
	err := v.vocabularyRepository.CreateOneByNameUnitInAdmin(ctx, vocabulary)

	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyUseCase) CreateOneInAdmin(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()
	err := v.vocabularyRepository.CreateOneInAdmin(ctx, vocabulary)

	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyUseCase) UpdateOneInAdmin(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	data, err := v.vocabularyRepository.UpdateOneInAdmin(ctx, vocabulary)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (v *vocabularyUseCase) UpdateOneImageInAdmin(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	data, err := v.vocabularyRepository.UpdateOneImageInAdmin(ctx, vocabulary)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (v *vocabularyUseCase) UpdateOneAudioInAdmin(c context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	ctx, cancel := context.WithTimeout(c, v.contextTimeout)
	defer cancel()

	err := v.vocabularyRepository.UpdateOneAudioInAdmin(ctx, vocabulary)
	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyUseCase) UpdateIsFavouriteInUser(ctx context.Context, vocabularyID string, isFavourite int) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	err := v.vocabularyRepository.UpdateIsFavouriteInUser(ctx, vocabularyID, isFavourite)
	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyUseCase) DeleteOneInAdmin(ctx context.Context, vocabularyID string) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	err := v.vocabularyRepository.DeleteOneInAdmin(ctx, vocabularyID)
	if err != nil {
		return err
	}

	return err
}
