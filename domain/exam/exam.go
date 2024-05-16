package exam_domain

import (
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
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
	Learner    string    `bson:"learner" json:"learner"`
}

type ExamResponse struct {
	ID       primitive.ObjectID   `bson:"_id" json:"_id"`
	LessonID lesson_domain.Lesson `bson:"lesson" json:"lesson"`
	UnitID   unit_domain.Unit     `bson:"unit" json:"unit"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Duration    string `bson:"duration" json:"duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
	Learner    string    `bson:"learner" json:"learner"`

	IsComplete    int   `bson:"is_complete" json:"is_complete"`
	CountQuestion int64 `bson:"count_question" json:"count_question"`
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
	FetchMany(ctx context.Context, page string) ([]ExamResponse, DetailResponse, error)
	FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]ExamResponse, DetailResponse, error)
	FetchOneByUnitID(ctx context.Context, unitID string) (ExamResponse, error)

	CreateOne(ctx context.Context, exam *Exam) error
	UpdateOne(ctx context.Context, exam *Exam) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examID string) error

	UpdateCompleted(ctx context.Context, exam *Exam) error
	Statistics(ctx context.Context) (Statistics, error)
}
