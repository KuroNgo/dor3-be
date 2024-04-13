package quiz_usecase

import (
	quiz_domain "clean-architecture/domain/quiz"
	"context"
	"time"
)

type quizUseCase struct {
	quizRepository quiz_domain.IQuizRepository
	contextTimeout time.Duration
}

func (q *quizUseCase) FetchManyByLessonID(ctx context.Context, unitID string) (quiz_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (q *quizUseCase) FetchManyByUnitID(ctx context.Context, unitID string) (quiz_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (q *quizUseCase) UpdateCompleted(ctx context.Context, quizID string, isComplete int) error {
	//TODO implement me
	panic("implement me")
}

func NewQuizUseCase(quizRepository quiz_domain.IQuizRepository, timeout time.Duration) quiz_domain.IQuizUseCase {
	return &quizUseCase{
		quizRepository: quizRepository,
		contextTimeout: timeout,
	}
}

func (q *quizUseCase) FetchMany(ctx context.Context) (quiz_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	quiz, err := q.quizRepository.FetchMany(ctx)
	if err != nil {
		return quiz_domain.Response{}, err
	}

	return quiz, err
}

func (q *quizUseCase) UpdateOne(ctx context.Context, quizID string, quiz quiz_domain.Quiz) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizRepository.UpdateOne(ctx, quizID, quiz)
	if err != nil {
		return err
	}

	return err
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

	return err
}
