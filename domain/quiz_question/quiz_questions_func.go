package quiz_question_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	QuizID       primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`

	Content       string   `bson:"content" json:"content"`
	Type          string   `bson:"type" json:"type"`
	Level         int      `bson:"level" json:"level"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
}

type IQuizQuestionUseCase interface {
	FetchManyInAdmin(ctx context.Context, page string) (Response, error)
	FetchByIDInAdmin(ctx context.Context, id string) (QuizQuestion, error)
	FetchManyByQuizIDInAdmin(ctx context.Context, quizID string) (Response, error)
	FetchOneByQuizIDInAdmin(ctx context.Context, quizID string) (QuizQuestion, error)

	UpdateOneInAdmin(ctx context.Context, quizQuestion *QuizQuestion) (*mongo.UpdateResult, error)
	CreateOneInAdmin(ctx context.Context, quizQuestion *QuizQuestion) error
	DeleteOneInAdmin(ctx context.Context, quizID string) error
}
