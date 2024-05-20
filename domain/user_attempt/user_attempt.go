package user_attempt_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUserAttempt = "user_attempt"
)

type UserProcess struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExamID     primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	QuizID     primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`

	Score         float32   `bson:"score" json:"score"`
	ProcessStatus int       `bson:"process_status" json:"process_status"`
	CompletedDate time.Time `bson:"completed_date" json:"completed_date"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type Response struct {
	Statistics  Statistics  `json:"statistics"`
	UserProcess UserProcess `bson:"user_process" json:"user_process"`
}

type Statistics struct {
	TotalScore   int64 `bson:"total_score" json:"total_score"`
	AverageScore int8  `bson:"average_score" json:"average_score"`
}

type IUserProcessRepository interface {
	FetchManyByUserID(ctx context.Context, userID string) (Response, error)
	FetchOneByUnitIDAndUserID(ctx context.Context, userID string, unit string) (UserProcess, error)
	CreateAttemptByExerciseID(ctx context.Context, userID UserProcess) error
	UpdateAttemptByUserID(ctx context.Context, userID UserProcess) error
	UpdateAttemptByExamID(ctx context.Context, userID UserProcess) error
	UpdateAttemptByQuizID(ctx context.Context, userID UserProcess) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}
