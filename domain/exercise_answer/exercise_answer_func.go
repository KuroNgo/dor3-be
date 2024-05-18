package exercise_answer_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Input struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Answer      string    `bson:"answer" json:"answer"`
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type IExerciseAnswerUseCase interface {
	FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOne(ctx context.Context, exerciseAnswer *ExerciseAnswer) error
	DeleteOne(ctx context.Context, exerciseID string) error
	DeleteAllAnswerByExerciseID(ctx context.Context, exerciseId string) error
}
