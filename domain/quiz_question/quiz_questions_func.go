package quiz_question

import (
	"context"
)

type Input struct {
	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`
	Level   int    `bson:"level" json:"level"`
}

type IQuizQuestionUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, quizID string) (Response, error)
	UpdateOne(ctx context.Context, quizQuestionID string, quizQuestion QuizQuestion) error
	CreateOne(ctx context.Context, quizQuestion *QuizQuestion) error
	DeleteOne(ctx context.Context, quizID string) error
}
