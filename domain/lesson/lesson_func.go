package lesson_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	CourseID primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name     string             `bson:"name" json:"name"`
	Content  string             `bson:"content" json:"content"`
	Image    string             `bson:"image" json:"image"`
	Level    int                `bson:"level" json:"level"`
}

type Update struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	CourseID primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name     string             `bson:"name" json:"name"`
	Content  string             `bson:"content" json:"content"`
	Image    string             `bson:"image" json:"image"`
	Level    int                `bson:"level" json:"level"`
}

//go:generate mockery --name ICourseUseCase
type ILessonUseCase interface {
	FetchMany(ctx context.Context, page string) ([]LessonResponse, DetailResponse, error)
	FetchManyNotPagination(ctx context.Context) ([]LessonResponse, DetailResponse, error)
	FetchByID(ctx context.Context, lessonID string) (LessonResponse, error)
	FindCourseIDByCourseName(ctx context.Context, courseName string) (primitive.ObjectID, error)
	FetchByIdCourse(ctx context.Context, idCourse string, page string) ([]LessonResponse, DetailResponse, error)

	CreateOne(ctx context.Context, lesson *Lesson) error
	CreateOneByNameCourse(ctx context.Context, lesson *Lesson) error

	DeleteOne(ctx context.Context, lessonID string) error

	UpdateImage(ctx context.Context, lesson *Lesson) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, lesson *Lesson) (*mongo.UpdateResult, error)
	// UpdateComplete automation
	UpdateComplete(ctx context.Context, lessonID string, lesson Lesson) error
}
