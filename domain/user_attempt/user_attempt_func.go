package user_attempt_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auto struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExamID     primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	QuizID     primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`

	Score         int `bson:"score" json:"score"`
	ProcessStatus int `bson:"process_status" json:"process_status"`
}

type IUserProcessUseCase interface {
	FetchManyByUserID(ctx context.Context, userID string) (Response, error)
	FetchOneByUnitID(ctx context.Context, unitID string) (UserProcess, error)
	CreateAttemptByExerciseID(ctx context.Context, userID UserProcess) error
	UpdateAttemptByUserID(ctx context.Context, userID UserProcess) error
	UpdateAttemptByExamID(ctx context.Context, userID UserProcess) error
	UpdateAttemptByQuizID(ctx context.Context, userID UserProcess) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}
