package exam_options_domain

import (
	exam_question_domain "clean-architecture/domain/exam_question"
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

	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type ExamOptionsResponse struct {
	ID       primitive.ObjectID                `bson:"_id" json:"_id"`
	Question exam_question_domain.ExamQuestion `bson:"question" json:"question"`

	Options       []string `bson:"options" json:"options"`
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
