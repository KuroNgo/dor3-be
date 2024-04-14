package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUserQuizAttempt = "user_quiz_process"
)

type UserQuizAttempt struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuizID        primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	Score         int64              `bson:"score" json:"score"`
	ProcessStatus int                `bson:"process_status" json:"process_status"`
	CompletedDate time.Time          `bson:"completed_date" json:"completed_date"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type Response struct {
	UserQuizAttempt            []UserQuizAttempt
	StudyTime                  time.Time `bson:"study_time" json:"study_time"`
	TotalNumberOfLessonLearned int16     `bson:"total_number_of_lesson_learned" json:"total_number_of_lesson_learned"`
	TotalScore                 int64     `bson:"total_score" json:"total_score"`
	AverageScore               int8      `bson:"average_score" json:"average_score"`
}

type IUserQuizRepository interface {
	FetchByIdLesson(ctx context.Context, idLesson string) (Response, error)
	CreateOne(ctx context.Context, userProcess *UserQuizAttempt) error
	DeleteOne(ctx context.Context, userProcessID string) error
}
