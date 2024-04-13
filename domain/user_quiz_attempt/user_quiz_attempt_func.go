package domain

import "context"

type IUserQuizUseCase interface {
	FetchByIdLesson(ctx context.Context, idLesson string) (Response, error)
	CreateOne(ctx context.Context, userProcess *UserQuizAttempt) error
	UpdateOne(ctx context.Context, userProcessID string, userProcess UserQuizAttempt) error
	DeleteOne(ctx context.Context, userProcessID string) error
}
