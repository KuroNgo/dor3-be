package quiz_result_usecase

import (
	quiz_result_domain "clean-architecture/domain/quiz_result"
	"golang.org/x/net/context"
	"time"
)

type quizResultUseCase struct {
	quizQuestionRepository quiz_result_domain.IQuizResultRepository
	contextTimeout         time.Duration
}

func NewQuizQuestionUseCase(quizQuestionRepository quiz_result_domain.IQuizResultRepository, timeout time.Duration) quiz_result_domain.IQuizResultUseCase {
	return &quizResultUseCase{
		quizQuestionRepository: quizQuestionRepository,
		contextTimeout:         timeout,
	}
}

func (q *quizResultUseCase) FetchManyInUser(ctx context.Context, page string) (quiz_result_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.FetchManyInUser(ctx, page)
	if err != nil {
		return quiz_result_domain.Response{}, err
	}

	return data, nil
}

func (q *quizResultUseCase) FetchManyByQuizIDInUser(ctx context.Context, quizID string) (quiz_result_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.FetchManyByQuizIDInUser(ctx, quizID)
	if err != nil {
		return quiz_result_domain.Response{}, err
	}

	return data, nil
}

func (q *quizResultUseCase) GetResultsByUserIDAndQuizIDInUser(ctx context.Context, userID string, quizID string) (quiz_result_domain.QuizResult, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.GetResultsByUserIDAndQuizIDInUser(ctx, userID, quizID)
	if err != nil {
		return quiz_result_domain.QuizResult{}, err
	}

	return data, nil
}

func (q *quizResultUseCase) CreateOneInUser(ctx context.Context, quizResult *quiz_result_domain.QuizResult) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizQuestionRepository.CreateOneInUser(ctx, quizResult)
	if err != nil {
		return err
	}

	return nil
}

func (q *quizResultUseCase) DeleteOneInUser(ctx context.Context, quizResultID string) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizQuestionRepository.DeleteOneInUser(ctx, quizResultID)
	if err != nil {
		return err
	}

	return nil
}
