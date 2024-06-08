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

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Duration    string `bson:"duration" json:"duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type DetailResponse struct {
	CountExercise int64      `bson:"count_exercise" json:"count_exercise"`
	Page          int64      `json:"page" bson:"page"`
	CurrentPage   int        `json:"current_page" bson:"current_page"`
	Statistics    Statistics `json:"statistics" bson:"statistics"`
}

type Statistics struct {
	Total int64 `bson:"total" json:"total"`
}

type IExerciseRepository interface {
	FetchMany(ctx context.Context, page string) ([]Exercise, DetailResponse, error)
	FetchByID(ctx context.Context, id string) (Exercise, error)
	FetchOneByUnitID(ctx context.Context, unitID string) (Exercise, error)
	FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]Exercise, DetailResponse, error)

	CreateOne(ctx context.Context, exercise *Exercise) error
	UpdateOne(ctx context.Context, exercise *Exercise) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, exerciseID string) error
	Statistics(ctx context.Context) (Statistics, error)
}
