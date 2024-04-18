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
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID   primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	Name       string             `bson:"name" json:"name"`
	ImageURL   string             `bson:"image_url" json:"image_url"`
	Level      int                `bson:"level" json:"level"`
	IsComplete int                `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdates string             `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	Unit []Unit
}

//go:generate mockery --name IUnitRepository
type IUnitRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FindLessonIDByLessonName(ctx context.Context, lessonName string) (primitive.ObjectID, error)
	FetchByIdLesson(ctx context.Context, idLesson string) (Response, error)

	CreateOne(ctx context.Context, unit *Unit) error
	CreateOneByNameLesson(ctx context.Context, unit *Unit) error

	UpdateOne(ctx context.Context, unit *Unit) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, unitID string) error

	// UpdateComplete automation
	UpdateComplete(ctx context.Context, update Update) error
	CheckLessonComplete(ctx context.Context, lessonID primitive.ObjectID) (bool, error)
}
