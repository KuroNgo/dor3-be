package quiz_answer

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionQuizAnswers = "quiz_answer"
)

type QuizAnswer struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content     string    `bson:"content" json:"content"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type Response struct {
	QuizAnswer []QuizAnswer
}

type IQuizAnswerRepository interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, quizAnswer *QuizAnswer) error
	DeleteOne(ctx context.Context, quizID string) error
}
