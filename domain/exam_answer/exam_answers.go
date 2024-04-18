package exam_answer_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionExamAnswers = "exam_answer"
)

type ExamAnswer struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content     string    `bson:"content" json:"content"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type Response struct {
	ExamAnswer []ExamAnswer
}

type IExamAnswerRepository interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, examAnswer *ExamAnswer) error
	DeleteOne(ctx context.Context, examID string) error
	DeleteAll(ctx context.Context) error
}
