package exercise_result_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Input struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`

	Score     int16     `bson:"score" json:"score"`
	StartedAt time.Time `bson:"started_at" json:"started_at"`
	Status    int       `bson:"status" json:"status"`
}

type IExerciseResultUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, userID string) (Response, error)

	GetResultsByUserIDAndExamID(ctx context.Context, userID string, exerciseID string) (ExerciseResult, error)

	CreateOne(ctx context.Context, exerciseResult *ExerciseResult) error
	UpdateStatus(ctx context.Context, exerciseResultID string, status int) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, exerciseResultID string) error

	CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int
	CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64
}
