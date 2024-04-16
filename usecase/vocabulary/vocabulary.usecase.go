package vocabulary_usecase

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type vocabularyUseCase struct {
	vocabularyRepository vocabulary_domain.IVocabularyRepository
	contextTimeout       time.Duration
}

func (v *vocabularyUseCase) GetLatestVocabulary(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.GetLatestVocabulary(ctx)
	if err != nil {
		return nil, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) GetLatestVocabularyBatch(ctx context.Context) (vocabulary_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func NewVocabularyUseCase(vocabularyRepository vocabulary_domain.IVocabularyRepository, timeout time.Duration) vocabulary_domain.IVocabularyUseCase {
	return &vocabularyUseCase{
		vocabularyRepository: vocabularyRepository,
		contextTimeout:       timeout,
	}
}

func (v *vocabularyUseCase) FetchByIdUnit(ctx context.Context, idUnit string) (vocabulary_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.FetchByIdUnit(ctx, idUnit)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) FetchByWord(ctx context.Context, word string) (vocabulary_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.FetchByWord(ctx, word)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) FetchByLesson(ctx context.Context, lessonName string) (vocabulary_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.FetchByLesson(ctx, lessonName)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) FetchMany(ctx context.Context, page string) (vocabulary_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.FetchMany(ctx, page)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	return vocabulary, err
}

func (v *vocabularyUseCase) UpdateOneAudio(c context.Context, vocabularyID string, linkURL string) error {
	ctx, cancel := context.WithTimeout(c, v.contextTimeout)
	defer cancel()

	err := v.vocabularyRepository.UpdateOneAudio(ctx, vocabularyID, linkURL)
	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyUseCase) UpdateIsFavourite(ctx context.Context, vocabularyID string, isFavourite int) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	err := v.vocabularyRepository.UpdateIsFavourite(ctx, vocabularyID, isFavourite)
	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyUseCase) GetAllVocabulary(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	vocabulary, err := v.vocabularyRepository.GetAllVocabulary(ctx)
	if err != nil {
		return nil, err
	}

	return vocabulary, nil
}

func (v *vocabularyUseCase) FindUnitIDByUnitLevel(ctx context.Context, unitLevel int) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	unitID, err := v.vocabularyRepository.FindUnitIDByUnitLevel(ctx, unitLevel)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return unitID, err
}

func (v *vocabularyUseCase) CreateOneByNameUnit(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()
	err := v.vocabularyRepository.CreateOne(ctx, vocabulary)

	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyUseCase) UpdateOne(ctx context.Context, vocabularyID string, vocabulary vocabulary_domain.Vocabulary) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	err := v.vocabularyRepository.UpdateOne(ctx, vocabularyID, vocabulary)
	if err != nil {
		return err
	}

	return err
}

func (v *vocabularyUseCase) CreateOne(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()
	err := v.vocabularyRepository.CreateOne(ctx, vocabulary)

	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyUseCase) DeleteOne(ctx context.Context, vocabularyID string) error {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	err := v.vocabularyRepository.DeleteOne(ctx, vocabularyID)
	if err != nil {
		return err
	}

	return err
}
