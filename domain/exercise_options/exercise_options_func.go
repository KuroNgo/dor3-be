package exercise_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	QuestionID    primitive.ObjectID `bson:"question_id" json:"question_id"`
	Content       string             `bson:"content" json:"content"`
	CorrectAnswer string             `bson:"correct_answer" json:"correct_answer"`
}

type IExerciseOptionUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, exerciseOptions *ExerciseOptions) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, exerciseOptions *ExerciseOptions) error
	DeleteOne(ctx context.Context, optionsID string) error
}
