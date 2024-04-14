package exam_question_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionExamQuestion = "exam_question"
)

type ExamQuestion struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	ExamID primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	ExamQuestion []ExamQuestion
}

type IExamQuestionRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string) (Response, error)
	UpdateOne(ctx context.Context, examQuestionID string, examQuestion ExamQuestion) error
	CreateOne(ctx context.Context, examQuestion *ExamQuestion) error
	DeleteOne(ctx context.Context, examID string) error
}
