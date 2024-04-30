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

	Score         int       `bson:"score" json:"score"`
	ProcessStatus int       `bson:"process_status" json:"process_status"`
	CompletedDate time.Time `bson:"completed_date" json:"completed_date"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}

type Response struct {
	UserProcess []UserProcess
}

type Statistics struct {
	StudyTime                  time.Time `bson:"study_time" json:"study_time"`
	TotalNumberOfLessonLearned int16     `bson:"total_number_of_lesson_learned" json:"total_number_of_lesson_learned"`
	TotalScore                 int64     `bson:"total_score" json:"total_score"`
	AverageScore               int8      `bson:"average_score" json:"average_score"`
}

type IUserProcessRepository interface {
	FetchManyByUserID(ctx context.Context) (Response, error)
	CreateOneByUserID(ctx context.Context, userID string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}
