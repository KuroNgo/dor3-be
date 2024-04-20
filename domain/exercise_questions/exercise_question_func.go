package exercise_questions_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	ExerciseID   primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`

	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`
	Level   int    `bson:"level" json:"level"`
}

type IExerciseQuestionUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExerciseID(ctx context.Context, exerciseID string) (Response, error)

	UpdateOne(ctx context.Context, exerciseQuestion *ExerciseQuestion) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, exerciseQuestion *ExerciseQuestion) error
	DeleteOne(ctx context.Context, exerciseID string) error
}
