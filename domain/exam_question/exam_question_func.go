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
	FetchManyInAdmin(ctx context.Context, page string) (Response, error)
	FetchQuestionByIDInAdmin(ctx context.Context, id string) (ExamQuestion, error)
	FetchManyByExamIDInAdmin(ctx context.Context, examID string, page string) (Response, error)
	FetchOneByExamIDInAdmin(ctx context.Context, examID string) (ExamQuestionResponse, error)

	CreateOneInAdmin(ctx context.Context, examQuestion *ExamQuestion) error
	UpdateOneInAdmin(ctx context.Context, examQuestion *ExamQuestion) (*mongo.UpdateResult, error)
	DeleteOneInAdmin(ctx context.Context, examID string) error
}
