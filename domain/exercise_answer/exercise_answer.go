package exercise_answer_domain

import (
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

type Response struct {
	ExerciseAnswer []ExerciseAnswer `json:"exercise_answer" bson:"exercise_answer"`
}

type IExerciseAnswerRepository interface {
	FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOne(ctx context.Context, exerciseAnswer *ExerciseAnswer) error
	DeleteOne(ctx context.Context, exerciseID string) error
	DeleteAllAnswerByExerciseID(ctx context.Context, exerciseId string) error
}
