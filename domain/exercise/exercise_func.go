package exercise_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Input struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`
}

type IExerciseUseCase interface {
	FetchMany(ctx context.Context, page string) ([]ExerciseResponse, int64, error)
	FetchManyByLessonID(ctx context.Context, unitID string) ([]ExerciseResponse, error)
	FetchManyByUnitID(ctx context.Context, unitID string) ([]ExerciseResponse, error)
	FetchOneByUnitID(ctx context.Context, unitID string) (ExerciseResponse, error)

	UpdateOne(ctx context.Context, exercise *Exercise) (*mongo.UpdateResult, error)
	UpdateCompleted(ctx context.Context, exerciseID string, isComplete int) error

	CreateOne(ctx context.Context, exercise *Exercise) error
	DeleteOne(ctx context.Context, exerciseID string) error
}
