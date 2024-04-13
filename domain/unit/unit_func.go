package unit_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	Name     string             `bson:"name" json:"name"`
	Content  string             `bson:"content" json:"content"`
}

type Update struct {
	UnitID     string `bson:"_id" json:"_id"`
	LessonID   string `bson:"lesson_id" json:"lesson_id"`
	IsComplete int    `bson:"is_complete" json:"is_complete"`
	WhoUpdate  string `bson:"who_update" json:"who_update"`
}

//go:generate mockery --name IUnitUseCase
type IUnitUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FindLessonIDByLessonName(ctx context.Context, lessonName string) (primitive.ObjectID, error)
	FetchByIdLesson(ctx context.Context, idLesson string) (Response, error)
	CreateOne(ctx context.Context, unit *Unit) error
	CreateOneByNameLesson(ctx context.Context, unit *Unit) error
	UpdateOne(ctx context.Context, unitID string, unit Unit) (Response, error)
	UpsertOne(ctx context.Context, id string, unit *Unit) (Response, error)
	DeleteOne(ctx context.Context, unitID string) error

	// UpdateComplete automation
	UpdateComplete(ctx context.Context, update Update) error
}
