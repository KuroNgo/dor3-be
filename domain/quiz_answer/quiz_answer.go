package quiz_answer_domain

import (
	quiz_question_domain "clean-architecture/domain/quiz_question"
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
	Learner     string    `bson:"learner" json:"learner"`
}

type QuizAnswerResponse struct {
	ID       primitive.ObjectID                `bson:"_id" json:"_id"`
	UserID   primitive.ObjectID                `bson:"user_id" json:"user_id"`
	Question quiz_question_domain.QuizQuestion `bson:"question" json:"question"`

	Answer      string    `bson:"answer" json:"answer"`
	IsCorrect   int       `bson:"correct" json:"correct"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type Response struct {
	QuizAnswer []QuizAnswerResponse `json:"quiz_answer" bson:"quiz_answer"`
	Score      int                  `json:"score" bson:"score"`
}

type IQuizAnswerRepository interface {
	FetchManyAnswerQuestionIDInUser(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOneInUser(ctx context.Context, quizAnswer *QuizAnswer) error
	DeleteOneInUser(ctx context.Context, quizID string) error
	DeleteAllAnswerByQuizIDInUser(ctx context.Context, quizId string) error
}
