package quiz_question_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	QuizID       primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`

	Content       string   `bson:"content" json:"content"`
	Type          string   `bson:"type" json:"type"`
	Level         int      `bson:"level" json:"level"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
}

type IQuizQuestionUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchByID(ctx context.Context, id string) (QuizQuestion, error)
	FetchManyByQuizID(ctx context.Context, quizID string) (Response, error)

	UpdateOne(ctx context.Context, quizQuestion *QuizQuestion) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, quizQuestion *QuizQuestion) error
	DeleteOne(ctx context.Context, quizID string) error
}
