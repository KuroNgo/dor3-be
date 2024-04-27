package quiz_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Input struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`
}

type Completed struct {
	LessonID   primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID     primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	IsComplete int                `bson:"is_complete" json:"is_complete"`
}

//go:generate mockery --name IQuizUseCase
type IQuizUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (Response, error)

	UpdateOne(ctx context.Context, quiz *Quiz) (*mongo.UpdateResult, error)
	UpdateCompleted(ctx context.Context, quiz *Quiz) error

	CreateOne(ctx context.Context, quiz *Quiz) error
	DeleteOne(ctx context.Context, quizID string) error
}
