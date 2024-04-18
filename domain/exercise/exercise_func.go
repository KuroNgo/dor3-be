package exercise_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Input struct {
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`
}

type IExerciseUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByLessonID(ctx context.Context, unitID string) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (Response, error)

	UpdateOne(ctx context.Context, exercise *Exercise) (*mongo.UpdateResult, error)
	UpdateCompleted(ctx context.Context, exerciseID string, isComplete int) error

	CreateOne(ctx context.Context, exercise *Exercise) error
	DeleteOne(ctx context.Context, exerciseID string) error
}
