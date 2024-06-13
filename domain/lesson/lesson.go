package lesson_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionLesson = "lesson"
)

type Lesson struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	CourseID primitive.ObjectID `bson:"course_id" json:"course_id"`

	Name     string `bson:"name" json:"name"`
	Content  string `bson:"content" json:"content"`
	ImageURL string `bson:"image_url" json:"image_url"`
	AssetURL string `bson:"asset_url" json:"asset_url"`
	Level    int    `bson:"level" json:"level"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type LessonResponse struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	CourseID primitive.ObjectID `bson:"course_id" json:"course_id"`

	Name     string `bson:"name" json:"name"`
	Content  string `bson:"content" json:"content"`
	ImageURL string `bson:"image_url" json:"image_url"`
	Level    int    `bson:"level" json:"level"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`

	CountVocabulary int32 `json:"count_vocabulary"`
	CountUnit       int32 `json:"count_unit"`
}

type DetailResponse struct {
	Page        int64      `json:"page"`
	CurrentPage int        `json:"current_page"`
	Statistics  Statistics `json:"statistics"`
}

type Statistics struct {
	Total           int64 `json:"total"`
	CountVocabulary int64 `json:"count_vocabulary"`
	CountUnit       int64 `json:"count_unit"`
}

//go:generate mockery --name ILessonRepository
type ILessonRepository interface {
	FetchManyNotPaginationInUser(ctx context.Context, userID primitive.ObjectID) ([]LessonProcessResponse, DetailResponse, error)
	FetchByIDInUser(ctx context.Context, userID primitive.ObjectID, lessonID string) (LessonProcessResponse, error)
	FetchByIDCourseInUser(ctx context.Context, userID primitive.ObjectID, courseID string, page string) ([]LessonProcessResponse, DetailResponse, error)
	FetchManyInUser(ctx context.Context, userID primitive.ObjectID, page string) ([]LessonProcessResponse, DetailResponse, error)
	UpdateCompleteInUser(ctx context.Context, user primitive.ObjectID) (*mongo.UpdateResult, error)

	FetchManyInAdmin(ctx context.Context, page string) ([]LessonResponse, DetailResponse, error)
	FetchManyNotPaginationInAdmin(ctx context.Context) ([]LessonResponse, DetailResponse, error)
	FetchByIDInAdmin(ctx context.Context, lessonID string) (LessonResponse, error)
	FetchByIdCourseInAdmin(ctx context.Context, idCourse string, page string) ([]LessonResponse, DetailResponse, error)
	FindLessonIDByLessonNameInAdmin(ctx context.Context, lessonName string) (primitive.ObjectID, error)

	CreateOneInAdmin(ctx context.Context, lesson *Lesson) error
	CreateOneByNameCourseInAdmin(ctx context.Context, lesson *Lesson) error
	DeleteOneInAdmin(ctx context.Context, lessonID string) error
	UpdateImageInAdmin(ctx context.Context, lesson *Lesson) (*mongo.UpdateResult, error)
	UpdateOneInAdmin(ctx context.Context, lesson *Lesson) (*mongo.UpdateResult, error)

	Statistics(ctx context.Context, countOptions bson.M) (Statistics, error)
}
