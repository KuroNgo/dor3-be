package exercise_questions

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	Content    string             `bson:"content" json:"content"`
	Type       string             `bson:"type" json:"type"`
}

type IExerciseQuestionUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, exerciseID string) (Response, error)
	UpdateOne(ctx context.Context, exerciseQuestionID string, exerciseQuestion ExerciseQuestion) error
	CreateOne(ctx context.Context, exerciseQuestion *ExerciseQuestion) error
	DeleteOne(ctx context.Context, exerciseID string) error
}
