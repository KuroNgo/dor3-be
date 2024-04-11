package user_exercise_attempt

import "context"

type IUserProcessUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	FetchByIdLesson(ctx context.Context, idLesson string) (Response, error)
	CreateOne(ctx context.Context, userProcess *UserExerciseProcess) error
	UpdateOne(ctx context.Context, userProcessID string, userProcess UserExerciseProcess) error
	UpsertOne(ctx context.Context, userProcessID string, userProcess *UserExerciseProcess) (Response, error)
	DeleteOne(ctx context.Context, userProcessID string) error
}
