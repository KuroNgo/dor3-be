package quiz_question_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	QuizID primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`
	Level   int    `bson:"level" json:"level"`
}

type IQuizQuestionUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, quizID string) (Response, error)

	UpdateOne(ctx context.Context, quizQuestion *QuizQuestion) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, quizQuestion *QuizQuestion) error
	DeleteOne(ctx context.Context, quizID string) error
}
