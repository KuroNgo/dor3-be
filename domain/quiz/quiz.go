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
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Duration    string `bson:"duration" json:"duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	CountExam   int64      `bson:"count_quiz" json:"count_quiz"`
	Page        int64      `bson:"page" json:"page"`
	CurrentPage int        `bson:"current_page" json:"current_page"`
	Statistics  Statistics `bson:"statistics" json:"statistics"`
}

type Statistics struct {
	Total int64 `bson:"total" json:"total"`
}

//go:generate mockery --name IQuizRepository
type IQuizRepository interface {
	FetchOneByUnitIDInUser(ctx context.Context, userID primitive.ObjectID, unitID string) (Quiz, error)
	UpdateCompletedInUser(ctx context.Context, quiz *Quiz) error

	FetchManyInAdmin(ctx context.Context, page string) ([]Quiz, Response, error)
	FetchByIDInAdmin(ctx context.Context, id string) (Quiz, error)
	FetchManyByUnitIDInAdmin(ctx context.Context, unitID string, page string) ([]Quiz, Response, error)
	FetchOneByUnitIDInAdmin(ctx context.Context, unitID string) (Quiz, error)

	CreateOneInAdmin(ctx context.Context, quiz *Quiz) error
	UpdateOneInAdmin(ctx context.Context, quiz *Quiz) (*mongo.UpdateResult, error)
	DeleteOneInAdmin(ctx context.Context, quizID string) error
	Statistics(ctx context.Context) (Statistics, error)
}
