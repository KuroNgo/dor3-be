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

	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	Duration    string    `bson:"duration" json:"duration"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	Learner     string    `bson:"learner" json:"learner"`
}

type Complete struct {
	Id         primitive.ObjectID `bson:"_id" json:"_id"`
	IsComplete int                `bson:"is_complete" json:"is_complete"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	Learner    string             `bson:"learner" json:"learner"`
}

type IExerciseUseCase interface {
	FetchManyInAdmin(ctx context.Context, page string) ([]Exercise, DetailResponse, error)
	FetchByIDInAdmin(ctx context.Context, id string) (Exercise, error)
	FetchOneByUnitIDInAdmin(ctx context.Context, unitID string) (Exercise, error)
	FetchManyByUnitIDInAdmin(ctx context.Context, unitID string, page string) ([]Exercise, DetailResponse, error)

	UpdateOneInAdmin(ctx context.Context, exercise *Exercise) (*mongo.UpdateResult, error)
	CreateOneInAdmin(ctx context.Context, exercise *Exercise) error
	DeleteOneInAdmin(ctx context.Context, exerciseID string) error
}
