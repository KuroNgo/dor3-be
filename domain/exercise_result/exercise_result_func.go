package exercise_result

import "context"

type IExerciseResultUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, exerciseID string) (Response, error)
	CreateOne(ctx context.Context, exerciseResult *ExerciseResult) error
	DeleteOne(ctx context.Context, exerciseResultID string) error
}
