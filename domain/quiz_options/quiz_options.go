package quiz_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionQuizOptions = "quiz_options"
)

type QuizOptions struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content       string `bson:"content" json:"content"`
	CorrectAnswer string `bson:"correct_answer" json:"correct_answer"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	QuizOptions []QuizOptions `json:"quiz_options" bson:"quiz_options"`
}

type IQuizOptionRepository interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, quizOptions *QuizOptions) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, quizOptions *QuizOptions) error
	DeleteOne(ctx context.Context, optionsID string) error
}
