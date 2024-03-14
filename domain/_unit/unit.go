package unit_domain

import (
	lesson_domain "clean-architecture/domain/lesson"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUnit = "_unit"
)

type Unit struct {
	ID         primitive.ObjectID `bson:"id" json:"id"`
	LessonID   primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	Name       string             `bson:"name" json:"name"`
	Content    string             `bson:"content" json:"content"`
	IsComplete int                `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdates time.Time          `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	ID         primitive.ObjectID `bson:"id" json:"id"`
	Lesson     lesson_domain.Lesson
	Name       string    `bson:"name" json:"name"`
	Content    string    `bson:"content" json:"content"`
	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates time.Time `bson:"who_updates" json:"who_updates"`
}

//go:generate mockery --name IUnitRepository
type IUnitRepository interface {
	FetchMany(ctx context.Context) ([]Response, error)
	CreateOne(ctx context.Context, unit *Unit) error
	UpdateOne(ctx context.Context, unitID string, unit Unit) error
	UpsertOne(ctx context.Context, id string, unit *Unit) (Response, error)
	DeleteOne(ctx context.Context, unitID string) error
}
