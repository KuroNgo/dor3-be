package course_domain

import (
	lesson_management_domain "clean-architecture/domain/user_process/lesson_management"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
}

//go:generate mockery --name ICourseUseCase
type ICourseUseCase interface {
	FetchManyInUser(ctx context.Context, userID primitive.ObjectID, page string) ([]lesson_management_domain.CourseProcess, DetailForManyResponse, error)
	FetchByIDInUser(ctx context.Context, userID primitive.ObjectID, courseID string) (lesson_management_domain.CourseProcess, error)
	UpdateCompleteInUser(ctx context.Context) (*mongo.UpdateResult, error)

	FetchManyForEachCourseInAdmin(ctx context.Context, page string) ([]CourseResponse, DetailForManyResponse, error)
	FetchByIDInAdmin(ctx context.Context, courseID string) (CourseResponse, error)
	FindCourseIDByCourseNameInAdmin(ctx context.Context, courseName string) (primitive.ObjectID, error)

	UpdateOneInAdmin(ctx context.Context, course *Course) (*mongo.UpdateResult, error)
	CreateOneInAdmin(ctx context.Context, course *Course) error
	DeleteOneInAdmin(ctx context.Context, courseID string) error
}
