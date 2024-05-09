package exam_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionExamOptions = "exam_options"
)

type ExamOptions struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Answer        []string `bson:"answer" json:"answer"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	ExamOptions []ExamOptions `json:"exam_options" bson:"exam_options"`
}

type IExamOptionRepository interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, examOptions *ExamOptions) error
	UpdateOne(ctx context.Context, examOptions *ExamOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, optionsID string) error
}
