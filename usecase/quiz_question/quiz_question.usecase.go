package quiz_question_usecase

import (
	quiz_question_domain "clean-architecture/domain/quiz_question"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"time"
)

type quizQuestionUseCase struct {
	quizQuestionRepository quiz_question_domain.IQuizQuestionRepository
	contextTimeout         time.Duration
}

func NewQuizQuestionUseCase(quizQuestionRepository quiz_question_domain.IQuizQuestionRepository, timeout time.Duration) quiz_question_domain.IQuizQuestionUseCase {
	return &quizQuestionUseCase{
		quizQuestionRepository: quizQuestionRepository,
		contextTimeout:         timeout,
	}
}

func (q *quizQuestionUseCase) FetchMany(ctx context.Context, page string) (quiz_question_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.FetchMany(ctx, page)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}

	return data, nil
}

func (q *quizQuestionUseCase) FetchManyByQuizID(ctx context.Context, quizID string) (quiz_question_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.FetchManyByQuizID(ctx, quizID)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}

	return data, nil
}

func (q *quizQuestionUseCase) UpdateOne(ctx context.Context, quizQuestion *quiz_question_domain.QuizQuestion) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.UpdateOne(ctx, quizQuestion)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (q *quizQuestionUseCase) CreateOne(ctx context.Context, quizQuestion *quiz_question_domain.QuizQuestion) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizQuestionRepository.CreateOne(ctx, quizQuestion)
	if err != nil {
		return err
	}

	return nil
}

func (q *quizQuestionUseCase) DeleteOne(ctx context.Context, quizID string) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizQuestionRepository.DeleteOne(ctx, quizID)
	if err != nil {
		return err
	}

	return nil
}
