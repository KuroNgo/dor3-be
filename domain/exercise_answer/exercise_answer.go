package exercise_answer_domain

import (
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionExerciseAnswers = "exercise_answer"
)

type ExerciseAnswer struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Answer    string `bson:"answer" json:"answer"`
	IsCorrect int    `bson:"is_correct" json:"is_correct"`

	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
	Learner     string    `bson:"learner" json:"learner"`
}

type ExerciseAnswerResponse struct {
	ID       primitive.ObjectID                         `bson:"_id" json:"_id"`
	UserID   primitive.ObjectID                         `bson:"user_id" json:"user_id"`
	Question exercise_questions_domain.ExerciseQuestion `bson:"question" json:"question"`

	Answer      string    `bson:"answer" json:"answer"`
	IsCorrect   int       `bson:"correct" json:"correct"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type Response struct {
	ExerciseAnswer []ExerciseAnswerResponse `json:"exercise_answer" bson:"exercise_answer"`
	Score          int                      `json:"score" bson:"score"`
}

type IExerciseAnswerRepository interface {
	FetchManyAnswerQuestionIDInUser(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOneInUser(ctx context.Context, exerciseAnswer *ExerciseAnswer) error
	DeleteOneInUser(ctx context.Context, exerciseID string) error
	DeleteAllAnswerByExerciseIDInUser(ctx context.Context, exerciseId string) error
}
