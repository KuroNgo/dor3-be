package exam_result_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionExamResult = "exam_result"
)

type ExamResult struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExamID primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Score      int16     `bson:"score" json:"score"`
	StartedAt  time.Time `bson:"started_at" json:"started_at"`
	IsComplete int       `bson:"is_complete" json:"is_complete"`
}

type Response struct {
	TotalScore   int16        `bson:"total_score" json:"total_score"`
	AverageScore float64      `bson:"average_score" json:"average_score"`
	Percentage   float64      `bson:"percentage" json:"percentage"`
	Page         int64        `bson:"page" json:"page"`
	ExamResult   []ExamResult `json:"exam_result" bson:"exam_result"`
}

type Statistics struct {
	TotalScore   int16   `bson:"total_score" json:"total_score"`
	AverageScore float64 `bson:"average_score" json:"average_score"`
	Percentage   float64 `bson:"percentage" json:"percentage"`
}

type IExamResultRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, userID string) (Response, error)
	FetchManyByUserID(ctx context.Context, examID string) (Response, error)

	GetResultsByUserIDAndExamID(ctx context.Context, userID string, examID string) (ExamResult, error)
	GetAverageScoreByUser(ctx context.Context, userID string) (float64, error)
	GetOverallPerformance(ctx context.Context, userID string) (float64, error)

	CreateOne(ctx context.Context, examResult *ExamResult) error
	UpdateStatus(ctx context.Context, examResultID string, status int) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examResultID string) error

	CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int
	CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64
}
