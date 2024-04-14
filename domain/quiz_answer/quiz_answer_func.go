package quiz_answer

import "context"

type IQuizAnswerUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, quizAnswer *QuizAnswer) error
	DeleteOne(ctx context.Context, quizID string) error
}
