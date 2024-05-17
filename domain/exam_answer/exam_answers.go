package exam_answer_domain

import (
	exam_question_domain "clean-architecture/domain/exam_question"
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

	Answer      string    `bson:"answer" json:"answer"`
	IsCorrect   int       `bson:"correct" json:"correct"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type ExamAnswerResponse struct {
	ID       primitive.ObjectID                `bson:"_id" json:"_id"`
	UserID   primitive.ObjectID                `bson:"user_id" json:"user_id"`
	Question exam_question_domain.ExamQuestion `bson:"question" json:"question"`

	Answer      string    `bson:"answer" json:"answer"`
	IsCorrect   int       `bson:"correct" json:"correct"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type Response struct {
	ExamAnswerResponse []ExamAnswerResponse `json:"exam_answer" bson:"exam_answer"`
}

type IExamAnswerRepository interface {
	FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOne(ctx context.Context, examAnswer *ExamAnswer) error
	DeleteOne(ctx context.Context, examID string) error
	DeleteAllAnswerByExamID(ctx context.Context, examID string) error
}
