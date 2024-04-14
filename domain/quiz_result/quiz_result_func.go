package quiz_result

import "context"

type IQuizResultUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByQuizID(ctx context.Context, quizID string) (Response, error)
	CreateOne(ctx context.Context, quizResult *QuizResult) error
	DeleteOne(ctx context.Context, quizResultID string) error
}
