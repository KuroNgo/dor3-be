package quiz_usecase

import (
	quiz_domain "clean-architecture/domain/quiz"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type quizUseCase struct {
	quizRepository quiz_domain.IQuizRepository
	contextTimeout time.Duration
}

func NewQuizUseCase(quizRepository quiz_domain.IQuizRepository, timeout time.Duration) quiz_domain.IQuizUseCase {
	return &quizUseCase{
		quizRepository: quizRepository,
		contextTimeout: timeout,
	}
}

func (q *quizUseCase) FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]quiz_domain.QuizResponse, quiz_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	quiz, detail, err := q.quizRepository.FetchManyByUnitID(ctx, unitID, page)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}

	return quiz, detail, nil
}

func (q *quizUseCase) FetchManyByLessonID(ctx context.Context, unitID string, page string) ([]quiz_domain.QuizResponse, quiz_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	quiz, detail, err := q.quizRepository.FetchManyByLessonID(ctx, unitID, page)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}

	return quiz, detail, nil
}

func (q *quizUseCase) FetchOneByUnitID(ctx context.Context, unitID string) (quiz_domain.QuizResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	quiz, err := q.quizRepository.FetchOneByUnitID(ctx, unitID)
	if err != nil {
		return quiz_domain.QuizResponse{}, err
	}

	return quiz, nil
}

func (q *quizUseCase) UpdateCompleted(ctx context.Context, quiz *quiz_domain.Quiz) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizRepository.UpdateCompleted(ctx, quiz)
	if err != nil {
		return err
	}

	return nil
}

func (q *quizUseCase) FetchMany(ctx context.Context, page string) ([]quiz_domain.QuizResponse, quiz_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	quiz, detail, err := q.quizRepository.FetchMany(ctx, page)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}

	return quiz, detail, nil
}

func (q *quizUseCase) UpdateOne(ctx context.Context, quiz *quiz_domain.Quiz) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizRepository.UpdateOne(ctx, quiz)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (q *quizUseCase) CreateOne(ctx context.Context, quiz *quiz_domain.Quiz) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()
	err := q.quizRepository.CreateOne(ctx, quiz)

	if err != nil {
		return err
	}

	return nil
}

func (q *quizUseCase) DeleteOne(ctx context.Context, quizID string) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizRepository.DeleteOne(ctx, quizID)
	if err != nil {
		return err
	}

	return nil
}
