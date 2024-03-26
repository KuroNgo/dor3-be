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
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	CourseID    primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name        string             `bson:"name" json:"name"`
	Content     string             `bson:"content" json:"content"`
	Image       string             `bson:"image" json:"image"`
	Level       int                `bson:"level" json:"level"`
	IsCompleted int                `bson:"is_completed" json:"is_completed"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdates  string             `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	Lesson []Lesson `json:"data" bson:"data"`
}

//go:generate mockery --name ILessonRepository
type ILessonRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	FetchByIdCourse(ctx context.Context, idCourse string) (Response, error)
	CreateOne(ctx context.Context, lesson *Lesson) error
	UpdateOne(ctx context.Context, lessonID string, lesson Lesson) error
	UpsertOne(ctx context.Context, id string, lesson *Lesson) (*Lesson, error)
	DeleteOne(ctx context.Context, lessonID string) error

	// UpdateComplete automation
	UpdateComplete(ctx context.Context, lessonID string, lesson Lesson) error
}
