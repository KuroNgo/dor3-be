package exercise_domain

import (
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
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

	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
	Learner    string    `bson:"learner" json:"learner"`
}

type ExerciseResponse struct {
	ID     primitive.ObjectID   `bson:"_id" json:"_id"`
	Lesson lesson_domain.Lesson `bson:"lesson_id" json:"lesson_id"`
	Unit   unit_domain.Unit     `bson:"unit_id" json:"unit_id"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Duration    string `bson:"duration" json:"duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
	Learner    string    `bson:"learner" json:"learner"`

	IsComplete    int   `bson:"is_complete" json:"is_complete"`
	CountQuestion int32 `bson:"count_question" json:"count_question"`
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
	FetchMany(ctx context.Context, page string) ([]ExerciseResponse, DetailResponse, error)
	FetchOneByUnitID(ctx context.Context, unitID string) (ExerciseResponse, error)
	FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]ExerciseResponse, DetailResponse, error)

	CreateOne(ctx context.Context, exercise *Exercise) error
	UpdateOne(ctx context.Context, exercise *Exercise) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, exerciseID string) error

	UpdateCompleted(ctx context.Context, exercise *Exercise) error
	Statistics(ctx context.Context) (Statistics, error)
}
