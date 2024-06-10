package course_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionCourse = "course"
)

type Course struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdated  string             `bson:"who_updated" json:"who_updated"`
}

type CourseResponse struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdated  string             `bson:"who_updated" json:"who_updated"`

	CountLesson     int32 `bson:"count_lesson" json:"count_lesson"`
	CountVocabulary int32 `bson:"count_vocabulary" json:"count_vocabulary"`
}

type DetailForManyResponse struct {
	Page        int64      `json:"page" bson:"page"`
	CurrentPage int        `json:"current_page" bson:"current_page"`
	CountCourse int64      `bson:"count_course" json:"count_course"`
	Statistics  Statistics `json:"statistics" bson:"statistics"`
}

type Statistics struct {
	Total int64 `bson:"total" json:"total"`
}

//go:generate mockery --name ICourseRepository
type ICourseRepository interface {
	FetchManyInUser(ctx context.Context, userID primitive.ObjectID, page string) ([]CourseProcess, DetailForManyResponse, error)
	FetchByIDInUser(ctx context.Context, userID primitive.ObjectID, courseID string) (CourseProcess, error)
	UpdateCompleteInUser(ctx context.Context) (*mongo.UpdateResult, error)

	FetchManyForEachCourseInAdmin(ctx context.Context, page string) ([]CourseResponse, DetailForManyResponse, error)
	FetchByIDInAdmin(ctx context.Context, courseID string) (CourseResponse, error)
	FindCourseIDByCourseNameInAdmin(ctx context.Context, courseName string) (primitive.ObjectID, error)

	CreateOneInAdmin(ctx context.Context, course *Course) error
	UpdateOneInAdmin(ctx context.Context, course *Course) (*mongo.UpdateResult, error)
	DeleteOneInAdmin(ctx context.Context, courseID string) error

	Statistics(ctx context.Context, countOptions bson.M) (Statistics, error)
}
