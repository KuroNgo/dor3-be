package user_exam_process_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUserExamProcess = "user_exam_process"
)

type UserExamProcess struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExamId        primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	Score         int                `bson:"score" json:"score"`
	ProcessStatus int                `bson:"process_status" json:"process_status"`
	CompletedDate time.Time          `bson:"completed_date" json:"completed_date"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type Response struct {
	UserProcess                []UserExamProcess `bson:"user_exam_process" json:"user_exam_process"`
	StudyTime                  time.Time         `bson:"study_time" json:"study_time"`
	TotalNumberOfLessonLearned int16             `bson:"total_number_of_lesson_learned" json:"total_number_of_lesson_learned"`
	TotalScore                 int64             `bson:"total_score" json:"total_score"`
	AverageScore               int8              `bson:"average_score" json:"average_score"`
}

type IUserExamProcessRepository interface {
	FetchByIdLesson(ctx context.Context, idLesson string) (Response, error)
	CreateOne(ctx context.Context, userProcess *UserExamProcess) error
	DeleteOne(ctx context.Context, userProcessID string) error
}
