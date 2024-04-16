package quiz_options_domain

import "context"

type IExamOptionUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, quizOptionsID string, quizOptions QuizOptions) error
	CreateOne(ctx context.Context, quizOptions *QuizOptions) error
	DeleteOne(ctx context.Context, optionsID string) error
}
