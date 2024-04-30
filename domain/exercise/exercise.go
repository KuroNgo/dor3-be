package exercise_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionExercise = "exercise"
)

type Exercise struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`

	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type ExerciseResponse struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`

	CountQuestion int32 `bson:"count_question" json:"count_question"`
}

type Response struct {
	CountExercise int64 `bson:"count_exercise" json:"count_exercise"`
	Page          int64 `bson:"page" json:"page"`
	Exercise      []Exercise
}

type Statistics struct {
	CountExercise int64 `bson:"count_exercise" json:"count_exercise"`
}

type IExerciseRepository interface {
	FetchMany(ctx context.Context, page string) ([]ExerciseResponse, int64, error)
	FetchManyByLessonID(ctx context.Context, unitID string) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (ExerciseResponse, error)

	CreateOne(ctx context.Context, exercise *Exercise) error

	UpdateOne(ctx context.Context, exercise *Exercise) (*mongo.UpdateResult, error)
	UpdateCompleted(ctx context.Context, exerciseID string, isComplete int) error

	DeleteOne(ctx context.Context, exerciseID string) error
}
