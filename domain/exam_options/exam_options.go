package exam_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionExamOptions = "exam_options"
)

type ExamOptions struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content string `bson:"content" json:"content"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	ExamOptions []ExamOptions
}

type IExamOptionRepository interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, examOptionsID string, examOptions ExamOptions) error
	CreateOne(ctx context.Context, examOptions *ExamOptions) error
	DeleteOne(ctx context.Context, optionsID string) error
}
