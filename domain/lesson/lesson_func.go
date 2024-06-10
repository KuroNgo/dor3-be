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
	FetchManyNotPaginationInUser(ctx context.Context, userID primitive.ObjectID) ([]LessonProcessResponse, DetailResponse, error)
	FetchByIDInUser(ctx context.Context, userID primitive.ObjectID, lessonID string) (LessonProcessResponse, error)
	FetchByIDCourseInUser(ctx context.Context, userID primitive.ObjectID, courseID string, page string) ([]LessonProcessResponse, DetailResponse, error)
	FetchManyInUser(ctx context.Context, userID primitive.ObjectID, page string) ([]LessonProcessResponse, DetailResponse, error)
	UpdateCompleteInUser(ctx context.Context) (*mongo.UpdateResult, error)

	FetchManyInAdmin(ctx context.Context, page string) ([]LessonResponse, DetailResponse, error)
	FetchManyNotPaginationInAdmin(ctx context.Context) ([]LessonResponse, DetailResponse, error)
	FetchByIDInAdmin(ctx context.Context, lessonID string) (LessonResponse, error)
	FindLessonIDByLessonNameInAdmin(ctx context.Context, lessonName string) (primitive.ObjectID, error)
	FetchByIdCourseInAdmin(ctx context.Context, idCourse string, page string) ([]LessonResponse, DetailResponse, error)

	CreateOneInAdmin(ctx context.Context, lesson *Lesson) error
	CreateOneByNameCourseInAdmin(ctx context.Context, lesson *Lesson) error
	DeleteOneInAdmin(ctx context.Context, lessonID string) error
	UpdateImageInAdmin(ctx context.Context, lesson *Lesson) (*mongo.UpdateResult, error)
	UpdateOneInAdmin(ctx context.Context, lesson *Lesson) (*mongo.UpdateResult, error)
}
