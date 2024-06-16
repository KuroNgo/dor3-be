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

func (q *quizQuestionUseCase) FetchManyInAdmin(ctx context.Context, page string) (quiz_question_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.FetchManyInAdmin(ctx, page)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}

	return data, nil
}

func (q *quizQuestionUseCase) FetchOneByQuizIDInAdmin(ctx context.Context, quizID string) (quiz_question_domain.QuizQuestion, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.FetchOneByQuizIDInAdmin(ctx, quizID)
	if err != nil {
		return quiz_question_domain.QuizQuestion{}, err
	}

	return data, nil
}

func (q *quizQuestionUseCase) FetchByIDInAdmin(ctx context.Context, id string) (quiz_question_domain.QuizQuestion, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.FetchByIDInAdmin(ctx, id)
	if err != nil {
		return quiz_question_domain.QuizQuestion{}, err
	}

	return data, nil
}

func (q *quizQuestionUseCase) FetchManyByQuizIDInAdmin(ctx context.Context, quizID string) (quiz_question_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.FetchManyByQuizIDInAdmin(ctx, quizID)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}

	return data, nil
}

func (q *quizQuestionUseCase) UpdateOneInAdmin(ctx context.Context, quizQuestion *quiz_question_domain.QuizQuestion) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	data, err := q.quizQuestionRepository.UpdateOneInAdmin(ctx, quizQuestion)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (q *quizQuestionUseCase) CreateOneInAdmin(ctx context.Context, quizQuestion *quiz_question_domain.QuizQuestion) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizQuestionRepository.CreateOneInAdmin(ctx, quizQuestion)
	if err != nil {
		return err
	}

	return nil
}

func (q *quizQuestionUseCase) DeleteOneInAdmin(ctx context.Context, quizID string) error {
	ctx, cancel := context.WithTimeout(ctx, q.contextTimeout)
	defer cancel()

	err := q.quizQuestionRepository.DeleteOneInAdmin(ctx, quizID)
	if err != nil {
		return err
	}

	return nil
}
