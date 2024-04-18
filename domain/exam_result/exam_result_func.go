package exam_result_domain

import (
	"context"
	"time"
)

type Input struct {
	Score     int16     `bson:"score" json:"score"`
	StartedAt time.Time `bson:"started_at" json:"started_at"`
	Status    int       `bson:"status" json:"status"`
}

type IExamResultUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string) (Response, error)

	GetResultsByUserIDAndExamID(ctx context.Context, userID string, examID string) (ExamResult, error)

	CreateOne(ctx context.Context, examResult *ExamResult) error
	UpdateStatus(ctx context.Context, examResultID string, status *int) error
	DeleteOne(ctx context.Context, examResultID string) error

	CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int
	CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64
}
