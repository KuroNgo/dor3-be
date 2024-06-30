package exam_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionExam = "exam"
)

type Exam struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Duration    string `bson:"duration" json:"duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type DetailResponse struct {
	CountExam   int64      `bson:"count_exam" json:"count_exam"`
	Page        int64      `json:"page" bson:"page"`
	CurrentPage int        `json:"current_page"`
	Statistics  Statistics `json:"statistics" bson:"statistics"`
}

type Statistics struct {
	Total int64 `bson:"total" json:"total"`
}

type IExamRepository interface {
	FetchOneByUnitIDInUser(ctx context.Context, userID primitive.ObjectID, unitID string) (Exam, error)
	UpdateCompletedInUser(ctx context.Context, exam *Exam) error

	FetchManyInAdmin(ctx context.Context, page string) ([]Exam, DetailResponse, error)
	FetchExamByIDInAdmin(ctx context.Context, id string) (Exam, error)
	FetchManyByUnitIDInAdmin(ctx context.Context, unitID string, page string) ([]Exam, DetailResponse, error)
	FetchOneByUnitIDInAdmin(ctx context.Context, unitID string) (Exam, error)

	CreateOneInAdmin(ctx context.Context, exam *Exam) error
	UpdateOneInAdmin(ctx context.Context, exam *Exam) (*mongo.UpdateResult, error)
	DeleteOneInAdmin(ctx context.Context, examID string) error
	Statistics(ctx context.Context) (Statistics, error)
}
