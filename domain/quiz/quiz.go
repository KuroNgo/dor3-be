package quiz

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionQuiz = "quiz"
)

type Quiz struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Question      string             `bson:"question"`
	Options       []string           `bson:"options"`
	CorrectAnswer string             `bson:"correct_answer"`
	QuestionType  string             `bson:"question_type"`
}

type Request struct {
	Question      string   `bson:"question"`
	Options       []string `bson:"options"`
	CorrectAnswer string   `bson:"correct_answer"`
	QuestionType  string   `bson:"question_type"`
}

type IQuiz interface {
	Fetch(ctx context.Context) ([]Quiz, error)
	Update(ctx context.Context, quizID string) (Quiz, error)
	Create(ctx context.Context, quiz *Quiz) error
	Delete()
}
