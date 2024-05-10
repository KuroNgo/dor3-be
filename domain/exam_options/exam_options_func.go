package exam_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	QuestionID    primitive.ObjectID `bson:"question_id" json:"question_id"`
	Options       []string           `bson:"options" json:"options"`
	CorrectAnswer string             `bson:"correct_answer" json:"correct_answer"`
}

type IExamOptionsUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, examOptions *ExamOptions) error
	UpdateOne(ctx context.Context, examOptions *ExamOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examID string) error
}
