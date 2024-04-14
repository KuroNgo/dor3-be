package quiz_domain

import (
	"context"
	"time"
)

type Input struct {
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`
}

//go:generate mockery --name IQuizUseCase
type IQuizUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (Response, error)
	UpdateOne(ctx context.Context, quizID string, quiz Quiz) error
	UpdateCompleted(ctx context.Context, quizID string, isComplete int) error
	CreateOne(ctx context.Context, quiz *Quiz) error
	DeleteOne(ctx context.Context, quizID string) error
}
