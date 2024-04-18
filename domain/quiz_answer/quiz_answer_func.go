package quiz_answer_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Auto struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content     string    `bson:"content" json:"content"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type IQuizAnswerUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, quizAnswer *QuizAnswer) error
	DeleteOne(ctx context.Context, quizID string) error
}
