package exam_question_domain

import (
	exam_options_domain "clean-architecture/domain/exam_options"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionExamQuestion = "exam_question"
)

type ExamQuestion struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	ExamID       primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`

	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`
	Level   int    `bson:"level" json:"level"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type ExamQuestionResponse struct {
	ID         primitive.ObjectID              `bson:"_id" json:"_id"`
	ExamID     primitive.ObjectID              `bson:"exam_id" json:"exam_id"`
	Vocabulary vocabulary_domain.Vocabulary    `bson:"vocabulary" json:"vocabulary"`
	Options    exam_options_domain.ExamOptions `bson:"options" json:"options"`

	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`
	Level   int    `bson:"level" json:"level"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	Page                 int64                  `bson:"page" json:"page"`
	CurrentPage          int                    `bson:"current_page" json:"current_page"`
	Statistics           Statistics             `bson:"statistics" json:"statistics"`
	ExamQuestionResponse []ExamQuestionResponse `json:"exam_question_response" bson:"exam_question_response"`
}

type Statistics struct {
	Count int64 `bson:"count" json:"count"`
}

type IExamQuestionRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string, page string) (Response, error)

	CreateOne(ctx context.Context, examQuestion *ExamQuestion) error
	UpdateOne(ctx context.Context, examQuestion *ExamQuestion) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examID string) error
}
