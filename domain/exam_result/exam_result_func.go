package exam_result_domain

import (
	"context"
)

type IExamOptionsUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string) (Response, error)
	UpdateOne(ctx context.Context, examResultID string, examResult ExamResult) error
	CreateOne(ctx context.Context, examResult *ExamResult) error
	DeleteOne(ctx context.Context, examResultID string) error
}
