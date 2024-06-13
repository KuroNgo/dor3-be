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

func (q *quizUseCase) FetchByIDInAdmin(ctx context.Context, id string) (quiz_domain.Quiz, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizRepository.FetchByIDInAdmin(ctx, id)
	if err != nil {
		return quiz_domain.Quiz{}, err
	}

	return data, nil
}

func NewQuizUseCase(quizRepository quiz_domain.IQuizRepository, timeout time.Duration) quiz_domain.IQuizUseCase {
	return &quizUseCase{
		quizRepository: quizRepository,
		contextTimeout: timeout,
	}
}

func (q *quizUseCase) FetchManyByUnitIDInAdmin(ctx context.Context, unitID string, page string) ([]quiz_domain.Quiz, quiz_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	quiz, detail, err := q.quizRepository.FetchManyByUnitIDInAdmin(ctx, unitID, page)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}

	return quiz, detail, nil
}

func (q *quizUseCase) FetchOneByUnitIDInAdmin(ctx context.Context, unitID string) (quiz_domain.Quiz, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	quiz, err := q.quizRepository.FetchOneByUnitIDInAdmin(ctx, unitID)
	if err != nil {
		return quiz_domain.Quiz{}, err
	}

	return quiz, nil
}

func (q *quizUseCase) FetchManyInAdmin(ctx context.Context, page string) ([]quiz_domain.Quiz, quiz_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	quiz, detail, err := q.quizRepository.FetchManyInAdmin(ctx, page)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}

	return quiz, detail, nil
}

func (q *quizUseCase) UpdateOneInAdmin(ctx context.Context, quiz *quiz_domain.Quiz) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizRepository.UpdateOneInAdmin(ctx, quiz)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (q *quizUseCase) CreateOneInAdmin(ctx context.Context, quiz *quiz_domain.Quiz) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()
	err := q.quizRepository.CreateOneInAdmin(ctx, quiz)

	if err != nil {
		return err
	}

	return nil
}

func (q *quizUseCase) DeleteOneInAdmin(ctx context.Context, quizID string) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizRepository.DeleteOneInAdmin(ctx, quizID)
	if err != nil {
		return err
	}

	return nil
}
