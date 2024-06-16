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

	Content       string   `bson:"content" json:"content"`
	Type          string   `bson:"type" json:"type"`
	Level         int      `bson:"level" json:"level"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
}

type IExerciseQuestionUseCase interface {
	FetchManyInAdmin(ctx context.Context, page string) (Response, error)
	FetchByIDInAdmin(ctx context.Context, id string) (ExerciseQuestionResponse, error)
	FetchManyByExerciseIDInAdmin(ctx context.Context, exerciseID string) (Response, error)
	FetchOneByExerciseIDInAdmin(ctx context.Context, exerciseID string) (ExerciseQuestionResponse, error)

	UpdateOneInAdmin(ctx context.Context, exerciseQuestion *ExerciseQuestion) (*mongo.UpdateResult, error)
	CreateOneInAdmin(ctx context.Context, exerciseQuestion *ExerciseQuestion) error
	DeleteOneInAdmin(ctx context.Context, exerciseID string) error
}
