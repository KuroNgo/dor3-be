package quiz_answer_domain

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

	Answer      string    `bson:"content" json:"content"`
	IsCorrect   int       `bson:"is_correct" json:"is_correct"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type Response struct {
	QuizAnswer []QuizAnswer `json:"quiz_answer" bson:"quiz_answer"`
}

type IQuizAnswerRepository interface {
	FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOne(ctx context.Context, quizAnswer *QuizAnswer) error
	DeleteOne(ctx context.Context, quizID string) error
	DeleteAllAnswerByQuizID(ctx context.Context, quizId string) error
}
