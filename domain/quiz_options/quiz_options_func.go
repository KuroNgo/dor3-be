package quiz_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content string `bson:"content" json:"content"`
}

type IExamOptionUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, quizOptions *QuizOptions) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, quizOptions *QuizOptions) error
	DeleteOne(ctx context.Context, optionsID string) error
}
