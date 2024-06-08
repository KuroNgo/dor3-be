package unit_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionUnit = "unit"
)

type Unit struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`

	Name     string `bson:"name" json:"name"`
	ImageURL string `bson:"image_url" json:"image_url"`
	Level    int    `bson:"level" json:"level"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	WhoCreate string    `bson:"who_create" json:"who_create"`
}

type UnitResponse struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`

	Name     string `bson:"name" json:"name"`
	ImageURL string `bson:"image_url" json:"image_url"`
	Level    int    `bson:"level" json:"level"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	WhoCreate string    `bson:"who_create" json:"who_create"`

	CountVocabulary int32 `json:"count_vocabulary"`
}

type DetailResponse struct {
	Page        int64 `json:"page"`
	CurrentPage int   `json:"current_page"`
}

type Statistics struct {
	CountUnit int64 `json:"count_unit"`
}

//go:generate mockery --name IUnitRepository
type IUnitRepository interface {
	FetchMany(ctx context.Context, page string) ([]UnitResponse, DetailResponse, error)
	FetchOneByID(ctx context.Context, id string) (UnitResponse, error)
	FetchManyNotPagination(ctx context.Context) ([]UnitResponse, error)
	FetchByIdLesson(ctx context.Context, idLesson string, page string) ([]UnitResponse, DetailResponse, error)

	CreateOne(ctx context.Context, unit *Unit) error
	CreateOneByNameLesson(ctx context.Context, unit *Unit) error
	DeleteOne(ctx context.Context, unitID string) error
	UpdateOne(ctx context.Context, unit *Unit) (*mongo.UpdateResult, error)
}
