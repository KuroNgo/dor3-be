package exam_options_domain

import (
	"context"
)

type Input struct {
	Content string `bson:"content" json:"content"`
}

type IExamOptionsUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, examOptionsID string, examOptions ExamOptions) error
	CreateOne(ctx context.Context, examOptions *ExamOptions) error
	DeleteOne(ctx context.Context, examID string) error
}
