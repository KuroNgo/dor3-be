package quiz_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Duration    string `bson:"duration" json:"duration"`
}

type Completed struct {
	LessonID   primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID     primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	IsComplete int                `bson:"is_complete" json:"is_complete"`
}

//go:generate mockery --name IQuizUseCase
type IQuizUseCase interface {
	FetchManyInAdmin(ctx context.Context, page string) ([]Quiz, Response, error)
	FetchByIDInAdmin(ctx context.Context, id string) (Quiz, error)
	FetchManyByUnitIDInAdmin(ctx context.Context, unitID string, page string) ([]Quiz, Response, error)
	FetchOneByUnitIDInAdmin(ctx context.Context, unitID string) (Quiz, error)

	UpdateOneInAdmin(ctx context.Context, quiz *Quiz) (*mongo.UpdateResult, error)
	CreateOneInAdmin(ctx context.Context, quiz *Quiz) error
	DeleteOneInAdmin(ctx context.Context, quizID string) error
}
