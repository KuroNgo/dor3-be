package exam_management

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUserExamManagement = "user_process"
)

type ExamManagement struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	QuizID     primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	ExamID     primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Score         float32   `bson:"score" json:"score"`
	ProcessStatus int       `bson:"process_status" json:"process_status"`
	CompletedDate time.Time `bson:"completed_date" json:"completed_date"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type Response struct {
	Statistics     Statistics     `json:"statistics"`
	ExamManagement ExamManagement `bson:"exam_management" json:"exam_management"`
}

type Statistics struct {
	TotalScore   int64 `bson:"total_score" json:"total_score"`
	AverageScore int8  `bson:"average_score" json:"average_score"`
}

type IUserProcessRepository interface {
	FetchManyByUserID(ctx context.Context, userID string) (Response, error)
	FetchOneByUnitIDAndUserID(ctx context.Context, userID string, unit string) (ExamManagement, error)
	CreateExamManagementByExerciseID(ctx context.Context, userID ExamManagement) error
	UpdateExamManagementByUserID(ctx context.Context, userID ExamManagement) error
	UpdateExamManagementByExamID(ctx context.Context, userID ExamManagement) error
	UpdateExamManagementByQuizID(ctx context.Context, userID ExamManagement) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}
