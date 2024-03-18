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

//go:generate mockery --name IUnitUseCase
type IUnitUseCase interface {
	FetchMany(ctx context.Context) ([]Response, error)
	CreateOne(ctx context.Context, unit *Unit) error
	UpdateOne(ctx context.Context, unitID string, unit Unit) (Response, error)
	UpsertOne(ctx context.Context, id string, unit *Unit) (Response, error)
	DeleteOne(ctx context.Context, unitID string) error
}
