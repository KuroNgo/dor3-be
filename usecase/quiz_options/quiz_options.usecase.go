package quiz_options

import (
	quiz_options_domain "clean-architecture/domain/quiz_options"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type quizOptionsUseCase struct {
	quizOptionsRepository quiz_options_domain.IQuizOptionRepository
	contextTimeout        time.Duration
}

func (q *quizOptionsUseCase) FetchManyByQuestionID(ctx context.Context, questionID string) (quiz_options_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizOptionsRepository.FetchManyByQuestionID(ctx, questionID)
	if err != nil {
		return quiz_options_domain.Response{}, err
	}

	return data, nil
}

func (q *quizOptionsUseCase) UpdateOne(ctx context.Context, quizOptions *quiz_options_domain.QuizOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizOptionsRepository.UpdateOne(ctx, quizOptions)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (q *quizOptionsUseCase) CreateOne(ctx context.Context, quizOptions *quiz_options_domain.QuizOptions) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizOptionsRepository.CreateOne(ctx, quizOptions)
	if err != nil {
		return err
	}

	return nil
}

func (q *quizOptionsUseCase) DeleteOne(ctx context.Context, optionsID string) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizOptionsRepository.DeleteOne(ctx, optionsID)
	if err != nil {
		return err
	}

	return nil
}

func NewQuizOptionsUseCase(quizOptionsRepository quiz_options_domain.IQuizOptionRepository, timeout time.Duration) quiz_options_domain.IQuizOptionRepository {
	return &quizOptionsUseCase{
		quizOptionsRepository: quizOptionsRepository,
		contextTimeout:        timeout,
	}
}
