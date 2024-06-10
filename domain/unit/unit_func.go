package unit_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	Name     string             `bson:"name" json:"name"`
	Level    int                `bson:"level" json:"level"`
}

type CompleteInput struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	IsComplete int                `bson:"is_complete" json:"is_complete"`
}

type Update struct {
	UnitID     primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID   primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	IsComplete int                `bson:"is_complete" json:"is_complete"`
	WhoUpdate  string             `bson:"who_update" json:"who_update"`
}

//go:generate mockery --name IUnitUseCase
type IUnitUseCase interface {
	UpdateCompleteInUser(ctx context.Context) (*mongo.UpdateResult, error)

	FetchManyInAdmin(ctx context.Context, page string) ([]UnitResponse, DetailResponse, error)
	FetchOneByIDInAdmin(ctx context.Context, id string) (UnitResponse, error)
	FetchManyNotPaginationInAdmin(ctx context.Context) ([]UnitResponse, error)
	FetchByIdLessonInAdmin(ctx context.Context, idLesson string, page string) ([]UnitResponse, DetailResponse, error)

	CreateOneInAdmin(ctx context.Context, unit *Unit) error
	CreateOneByNameLessonInAdmin(ctx context.Context, unit *Unit) error
	UpdateOneInAdmin(ctx context.Context, unit *Unit) (*mongo.UpdateResult, error)
	DeleteOneInAdmin(ctx context.Context, unitID string) error
}
