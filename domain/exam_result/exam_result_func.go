package exam_result_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Input struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExamID primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Score     int16     `bson:"score" json:"score"`
	StartedAt time.Time `bson:"started_at" json:"started_at"`
	Status    int       `bson:"is_complete" json:"is_complete"`
}

type IExamResultUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string) (Response, error)

	GetResultsByUserIDAndExamID(ctx context.Context, userID string, examID string) (ExamResult, error)
	GetResultByID(ctx context.Context, userID string) (ExamResult, error)

	CreateOne(ctx context.Context, examResult *ExamResult) error
	UpdateStatus(ctx context.Context, examResultID string, status int) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examResultID string) error

	CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int
	CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64
}
