package quiz_answer_usecase

import (
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	"context"
	"time"
)

type quizResultUseCase struct {
	quizAnswerRepository quiz_answer_domain.IQuizAnswerRepository
	contextTimeout       time.Duration
}

func (q *quizResultUseCase) FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (quiz_answer_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizAnswerRepository.FetchManyAnswerByUserIDAndQuestionID(ctx, questionID, userID)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}

	return data, nil
}

func (q *quizResultUseCase) CreateOne(ctx context.Context, quizAnswer *quiz_answer_domain.QuizAnswer) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizAnswerRepository.CreateOne(ctx, quizAnswer)
	if err != nil {
		return err
	}

	return nil
}

func (q *quizResultUseCase) DeleteOne(ctx context.Context, quizID string) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizAnswerRepository.DeleteOne(ctx, quizID)
	if err != nil {
		return err
	}

	return nil
}

func NewQuizResultUseCase(quizAnswerRepository quiz_answer_domain.IQuizAnswerRepository, timeout time.Duration) quiz_answer_domain.IQuizAnswerUseCase {
	return &quizResultUseCase{
		quizAnswerRepository: quizAnswerRepository,
		contextTimeout:       timeout,
	}
}
