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
	FetchManyInUser2(ctx context.Context, examID string, userID primitive.ObjectID) (Response, error)
	FetchManyInUser(ctx context.Context, page string) (Response, error)
	FetchManyByExerciseIDInUser(ctx context.Context, userID string) (Response, error)
	GetResultsExerciseIDInUser(ctx context.Context, userID string, exerciseID string) (ExerciseResult, error)

	CreateOneInUser(ctx context.Context, exerciseResult *ExerciseResult) error
	UpdateStatusInUser(ctx context.Context, exerciseResultID string, status int) (*mongo.UpdateResult, error)
	DeleteOneInUser(ctx context.Context, exerciseResultID string) error

	CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int
	CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64
}
