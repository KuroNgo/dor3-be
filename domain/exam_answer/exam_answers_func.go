package exam_answer_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Input struct {
	QuestionID  primitive.ObjectID `bson:"question_id" json:"question_id"`
	Content     string             `bson:"content" json:"content"`
	SubmittedAt time.Time          `bson:"submitted_at" json:"submitted_at"`
}

type IExamAnswerUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, examAnswerID string, examAnswer ExamAnswer) error
	CreateOne(ctx context.Context, examAnswer *ExamAnswer) error
	UpdateCompleted(ctx context.Context, examID string, isComplete int) error
	DeleteOne(ctx context.Context, examID string) error
}
