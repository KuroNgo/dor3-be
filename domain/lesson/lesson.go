package lesson_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	ImageURL    string             `bson:"image_url" json:"image_url"`
	Level       int                `bson:"level" json:"level"`
	IsCompleted int                `bson:"is_completed" json:"is_completed"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdates  string             `bson:"who_updates" json:"who_updates"`
}

type LessonResponse struct {
	ID              primitive.ObjectID `bson:"_id" json:"_id"`
	CourseID        primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name            string             `bson:"name" json:"name"`
	Content         string             `bson:"content" json:"content"`
	ImageURL        string             `bson:"image_url" json:"image_url"`
	Level           int                `bson:"level" json:"level"`
	IsCompleted     int                `bson:"is_completed" json:"is_completed"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdates      string             `bson:"who_updates" json:"who_updates"`
	CountVocabulary int32              `json:"count_vocabulary"`
	CountUnit       int32              `json:"count_unit"`
}

type Response struct {
	CountVocabulary int64 `json:"count_vocabulary"`
	Page            int64 `json:"page"`
	CountUnit       int64 `json:"count_unit"`
	Lesson          []Lesson
}

type Statistics struct {
	CountVocabulary int64 `json:"count_vocabulary"`
	CountUnit       int64 `json:"count_unit"`
}

//go:generate mockery --name ILessonRepository
type ILessonRepository interface {
	FetchMany(ctx context.Context) ([]LessonResponse, error)
	FindCourseIDByCourseName(ctx context.Context, courseName string) (primitive.ObjectID, error)
	FetchByIdCourse(ctx context.Context, idCourse string) (Response, error)

	CreateOne(ctx context.Context, lesson *Lesson) error
	CreateOneByNameCourse(ctx context.Context, lesson *Lesson) error

	UpdateOne(ctx context.Context, lesson *Lesson) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, lessonID string) error

	// UpdateComplete automation
	UpdateComplete(ctx context.Context, lessonID string, lesson Lesson) error
}
