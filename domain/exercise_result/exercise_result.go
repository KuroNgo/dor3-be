package exercise_result_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionExerciseResult = "exercise_result"
)

type ExerciseResult struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`

	Score      int16     `bson:"score" json:"score"`
	StartedAt  time.Time `bson:"started_at" json:"started_at"`
	IsComplete int       `bson:"is_complete" json:"is_complete"`
}

type Response struct {
	TotalScore     int16   `bson:"total_score" json:"total_score"`
	AverageScore   float64 `bson:"average_score" json:"average_score"`
	Percentage     float64 `bson:"percentage" json:"percentage"`
	Page           int64   `bson:"page" json:"page"`
	ExerciseResult []ExerciseResult
}

type IExerciseResultRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExerciseID(ctx context.Context, userID string) (Response, error)
	FetchManyByUserID(ctx context.Context, exerciseID string) (Response, error)

	GetResultsByUserIDAndExerciseID(ctx context.Context, userID string, exerciseID string) (ExerciseResult, error)
	GetAverageScoreByUser(ctx context.Context, userID string) (float64, error)
	GetOverallPerformance(ctx context.Context, userID string) (float64, error)

	CreateOne(ctx context.Context, exerciseResult *ExerciseResult) error
	UpdateStatus(ctx context.Context, exerciseResultID string, status int) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, exerciseResultID string) error

	CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int
	CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64
}
