package unit_domain

import (
	"context"
)

type Input struct {
	LessonName string `bson:"lesson_name" json:"lesson_name"`
	Name       string `bson:"name" json:"name"`
	Content    string `bson:"content" json:"content"`
}

//go:generate mockery --name IUnitUseCase
type IUnitUseCase interface {
	FetchMany(ctx context.Context) ([]Unit, error)
	FetchToDeleteMany(ctx context.Context) (*[]Unit, error)
	CreateOne(ctx context.Context, unit *Unit) error
	UpdateOne(ctx context.Context, unitID string, unit Unit) error
	UpsertOne(ctx context.Context, id string, unit *Unit) (*Unit, error)
	DeleteOne(ctx context.Context, unitID string) error
}
