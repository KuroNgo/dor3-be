package exam_question_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`
	Level   int    `bson:"level" json:"level"`
}

type IExamQuestionUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string) (Response, error)

	CreateOne(ctx context.Context, examQuestion *ExamQuestion) error
	UpdateOne(ctx context.Context, examQuestion *ExamQuestion) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examID string) error
}
