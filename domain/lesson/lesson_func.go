package lesson_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	CourseID primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name     string             `bson:"name" json:"name"`
	Content  string             `bson:"content" json:"content"`
}

//go:generate mockery --name ICourseUseCase
type ILessonUseCase interface {
	// FetchByID this func not alone api
	FetchByID(ctx context.Context, lessonID string) (*Lesson, error)

	FetchMany(ctx context.Context) ([]Lesson, error)
	FetchToDeleteMany(ctx context.Context) (*[]Lesson, error)
	UpdateOne(ctx context.Context, lessonID string, lesson Lesson) error
	CreateOne(ctx context.Context, lesson *Lesson) error
	UpsertOne(ctx context.Context, id string, lesson *Lesson) (*Lesson, error)
	DeleteOne(ctx context.Context, lessonID string) error
}
