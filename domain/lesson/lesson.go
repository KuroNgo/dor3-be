package lesson_domain

import (
	course_domain "clean-architecture/domain/course"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionLesson = "lesson"
)

type Lesson struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	CourseID   primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name       string             `bson:"name" json:"name"`
	Content    string             `bson:"content" json:"content"`
	Level      int                `bson:"level" json:"level"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdates string             `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Course     course_domain.Course
	Name       string    `bson:"name" json:"name"`
	Content    string    `bson:"content" json:"content"`
	Level      int       `bson:"level" json:"level"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

//go:generate mockery --name ILessonRepository
type ILessonRepository interface {
	FetchMany(ctx context.Context) ([]Response, error)
	CreateOne(ctx context.Context, lesson *Lesson) error
	UpdateOne(ctx context.Context, lessonID string, lesson Lesson) error
	UpsertOne(ctx context.Context, id string, lesson *Lesson) (*Lesson, error)
	DeleteOne(ctx context.Context, lessonID string) error
}
