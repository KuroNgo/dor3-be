package exam_domain

import (
	"context"
	"time"
)

type Input struct {
	Title       string
	Description string
	Duration    time.Duration
}

type IExamUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (Response, error)

	CreateOne(ctx context.Context, exam *Exam) error
	UpdateOne(ctx context.Context, examID string, exam Exam) error
	UpdateCompleted(ctx context.Context, examID string, isComplete int) error
	DeleteOne(ctx context.Context, examID string) error
}
