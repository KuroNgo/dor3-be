package exam_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	Content string `bson:"content" json:"content"`
}

type IExamOptionsUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, examOptions *ExamOptions) error
	UpdateOne(ctx context.Context, examOptions *ExamOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examID string) error
}
