package exercise_answer_domain

import (
	exercise_options_domain "clean-architecture/domain/exercise_options"
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

	Answer     string `bson:"answer" json:"answer"`
	BlankIndex int    `bson:"blank_index" json:"blank_index"`
	IsCorrect  int    `bson:"is_correct" json:"is_correct"`

	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type ExerciseAnswerResponse struct {
	ID       primitive.ObjectID                         `bson:"_id" json:"_id"`
	UserID   primitive.ObjectID                         `bson:"user_id" json:"user_id"`
	Question exercise_questions_domain.ExerciseQuestion `bson:"question" json:"question"`
	Options  exercise_options_domain.ExerciseOptions    `bson:"options" json:"options"`

	Answer     string `bson:"answer" json:"answer"`
	BlankIndex int    `bson:"blank_index" json:"blank_index"`
	IsCorrect  int    `bson:"is_correct" json:"is_correct"`

	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type Response struct {
	ExerciseAnswerResponse []ExerciseAnswerResponse `json:"exercise_answer" bson:"exercise_answer"`
}

type IExerciseAnswerRepository interface {
	FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOne(ctx context.Context, exerciseAnswer *ExerciseAnswer) error
	DeleteOne(ctx context.Context, exerciseID string) error
	DeleteAllAnswerByExerciseID(ctx context.Context, exerciseId string) error
}
