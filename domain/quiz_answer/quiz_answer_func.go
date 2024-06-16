package quiz_answer_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Input struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Answer      string    `bson:"content" json:"content"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type IQuizAnswerUseCase interface {
	FetchManyAnswerQuestionIDInUser(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOneInUser(ctx context.Context, quizAnswer *QuizAnswer) error
	DeleteOneInUser(ctx context.Context, quizID string) error
	DeleteAllAnswerByQuizIDInUser(ctx context.Context, quizId string) error
}
