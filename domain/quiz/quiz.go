package quiz_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionQuiz = "quiz"
)

type Quiz struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Duration    string `bson:"duration" json:"duration"`

	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type QuizResponse struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Duration    string `bson:"duration" json:"duration"`

	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`

	CountQuestion int64 `bson:"count_question" json:"count_question"`
}

type Response struct {
	Total       int64 `bson:"total" json:"total"`
	Page        int64 `bson:"page" json:"page"`
	CurrentPage int   `bson:"current_page" json:"current_page"`
}

type Statistics struct {
	Total int64 `bson:"total" json:"total"`
}

//go:generate mockery --name IQuizRepository
type IQuizRepository interface {
	FetchMany(ctx context.Context, page string) ([]QuizResponse, Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]QuizResponse, Response, error)
	FetchManyByLessonID(ctx context.Context, unitID string, page string) ([]QuizResponse, Response, error)
	FetchOneByUnitID(ctx context.Context, unitID string) (QuizResponse, error)

	CreateOne(ctx context.Context, quiz *Quiz) error

	UpdateOne(ctx context.Context, quiz *Quiz) (*mongo.UpdateResult, error)
	UpdateCompleted(ctx context.Context, quiz *Quiz) error

	DeleteOne(ctx context.Context, quizID string) error
}
