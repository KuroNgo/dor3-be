package exam_result_domain

import (
	"context"
)

type IExamResultUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string) (Response, error)
	CreateOne(ctx context.Context, examResult *ExamResult) error
	DeleteOne(ctx context.Context, examResultID string) error
}
