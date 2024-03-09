package lesson_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionLesson = "lesson"
)

type Lesson struct {
	ID         primitive.ObjectID `bson:"id" json:"id"`
	CourseID   primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name       string             `bson:"name" json:"name"`
	Content    string             `bson:"content" json:"content"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdates time.Time          `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	ID         primitive.ObjectID `bson:"id" json:"id"`
	CourseID   primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name       string             `bson:"name" json:"name"`
	Content    string             `bson:"content" json:"content"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdates time.Time          `bson:"who_updates" json:"who_updates"`
}

//go:generate mockery --name ILessonRepository
type ILessonRepository interface {
	FetchByID(ctx context.Context, lessonID string) (*Lesson, error)
	FetchMany(ctx context.Context) ([]Lesson, error)
	FetchToDeleteMany(ctx context.Context) (*[]Lesson, error)
	UpdateOne(ctx context.Context, lessonID string, lesson Lesson) error
	CreateOne(ctx context.Context, lesson *Lesson) error
	UpsertOne(ctx context.Context, id string, lesson *Lesson) (*Lesson, error)
	DeleteOne(ctx context.Context, lessonID string) error
}
