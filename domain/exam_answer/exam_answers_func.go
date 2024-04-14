package exam_answer_domain

import (
	"context"
)

type Input struct {
	Content string `bson:"content" json:"content"`
}

type IExamAnswerUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, examAnswer *ExamAnswer) error
	DeleteOne(ctx context.Context, examID string) error
}
