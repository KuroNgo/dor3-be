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
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LessonID     primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID       primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary" json:"vocabulary"`

	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`

	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	Count int64 `bson:"count" json:"count"`
	Quiz  []Quiz
}

//go:generate mockery --name IQuizRepository
type IQuizRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (Response, error)
	FetchManyByLessonID(ctx context.Context, unitID string) (Response, error)

	CreateOne(ctx context.Context, quiz *Quiz) error

	UpdateOne(ctx context.Context, quiz *Quiz) (*mongo.UpdateResult, error)
	UpdateCompleted(ctx context.Context, quiz *Quiz) error

	DeleteOne(ctx context.Context, quizID string) error
}
