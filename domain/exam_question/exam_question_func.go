package exam_question_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	ExamID  primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	Content string             `bson:"content" json:"content"`
	Type    string             `bson:"type" json:"type"`
}

type IExamQuestionUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string) (Response, error)
	UpdateOne(ctx context.Context, examQuestionID string, examQuestion ExamQuestion) error
	CreateOne(ctx context.Context, examQuestion *ExamQuestion) error
	DeleteOne(ctx context.Context, examID string) error
}
