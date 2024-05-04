package quiz_question_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionQuizQuestion = "quiz_question"
)

type QuizQuestion struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	QuizID       primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`

	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`
	Level   int    `bson:"level" json:"level"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	Page         int64          `bson:"page" json:"page"`
	QuizQuestion []QuizQuestion `json:"quiz_question" bson:"quiz_question"`
}

type IQuizQuestionRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByQuizID(ctx context.Context, quizID string) (Response, error)

	UpdateOne(ctx context.Context, quizQuestion *QuizQuestion) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, quizQuestion *QuizQuestion) error
	DeleteOne(ctx context.Context, quizID string) error
}
