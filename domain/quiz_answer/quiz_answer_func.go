package quiz_answer_domain

import "context"

type IQuizAnswerUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, quizAnswer *QuizAnswer) error
	DeleteOne(ctx context.Context, quizID string) error
}
