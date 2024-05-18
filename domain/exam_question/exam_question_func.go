package exam_question_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	ExamID       primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`

	Content       string   `bson:"content" json:"content"`
	Type          string   `bson:"type" json:"type"`
	Level         int      `bson:"level" json:"level"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
}

type IExamQuestionUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchQuestionByID(ctx context.Context, id string) (ExamQuestion, error)
	FetchManyByExamID(ctx context.Context, examID string, page string) (Response, error)

	CreateOne(ctx context.Context, examQuestion *ExamQuestion) error
	UpdateOne(ctx context.Context, examQuestion *ExamQuestion) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examID string) error
}
