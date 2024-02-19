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

func NewQuizUseCase(quizRepository quiz_domain.IQuizRepository, timeout time.Duration) quiz_domain.IQuizUseCase {
	return &quizUseCase{
		quizRepository: quizRepository,
		contextTimeout: timeout,
	}
}

func (q *quizUseCase) Fetch(ctx context.Context) ([]quiz_domain.Quiz, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	user, err := q.quizRepository.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (q *quizUseCase) FetchToDelete(ctx context.Context) (*[]quiz_domain.Quiz, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	user, err := q.quizRepository.FetchToDelete(ctx)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (q *quizUseCase) Update(ctx context.Context, quizID string, quiz quiz_domain.Quiz) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizRepository.Update(ctx, quizID, quiz)
	if err != nil {
		return err
	}

	return err
}

func (q *quizUseCase) Create(ctx context.Context, quiz *quiz_domain.Input) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()
	err := q.quizRepository.Create(ctx, quiz)

	if err != nil {
		return err
	}

	return nil
}

func (q *quizUseCase) Upsert(c context.Context, question string, quiz *quiz_domain.Quiz) (*quiz_domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, q.contextTimeout)
	defer cancel()
	return q.quizRepository.Upsert(ctx, question, quiz)
}

func (q *quizUseCase) Delete(ctx context.Context, quizID string) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizRepository.Delete(ctx, quizID)
	if err != nil {
		return err
	}

	return err
}
