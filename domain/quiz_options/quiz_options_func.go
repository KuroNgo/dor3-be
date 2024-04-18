package quiz_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type IExamOptionUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, quizOptions *QuizOptions) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, quizOptions *QuizOptions) error
	DeleteOne(ctx context.Context, optionsID string) error
}
