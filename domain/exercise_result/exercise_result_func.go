package exercise_result_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	FetchManyByExamID(ctx context.Context, exerciseID string) (Response, error)
	CreateOne(ctx context.Context, exerciseResult *ExerciseResult) error
	DeleteOne(ctx context.Context, exerciseResultID string) error
}
