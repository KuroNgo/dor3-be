package quiz_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionQuiz = "quiz"
)

type Quiz struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Question      string             `bson:"question" json:"question"`
	Options       []string           `bson:"options" json:"options"`
	CorrectAnswer string             `bson:"correct_answer" json:"correct_answer"`
	QuestionType  string             `bson:"question_type" json:"question_type"`
}

type Request struct {
	Question      string   `bson:"question,omitempty"`
	Options       []string `bson:"options,omitempty"`
	CorrectAnswer string   `bson:"correct_answer,omitempty"`

	// QuestionType can be included checkbox, check radius or write correct answer
	QuestionType string `bson:"question_type,omitempty"`
}

//go:generate mockery --name IQuizRepository
type IQuizRepository interface {
	Fetch(ctx context.Context) ([]Quiz, error)
	FetchToDelete(ctx context.Context) (*[]Quiz, error)
	Update(ctx context.Context, quizID string, quiz Quiz) error
	Create(ctx context.Context, quiz *Quiz) error
	Delete(ctx context.Context, quizID string) error
}
