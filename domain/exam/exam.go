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

	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type ExamResponse struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`

	CountQuestion int64 `bson:"count_question" json:"count_question"`
}

type Response struct {
	CountExam int64 `bson:"count_exam" json:"count_exam"`
	Page      int64 `json:"page" bson:"page"`
	Exam      []Exam
}

type IExamRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (Response, error)

	CreateOne(ctx context.Context, exam *Exam) error
	UpdateOne(ctx context.Context, exam *Exam) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examID string) error

	CountQuestion(ctx context.Context, examID string) int64
}
