package user_exam_process_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserProcessAuto struct {
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	CourseID      primitive.ObjectID `bson:"course_id" json:"course_id"`
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	Score         int                `bson:"score" json:"score"`
	ProcessStatus int                `bson:"process_status" json:"process_status"`
	CompletedDate time.Time          `bson:"completed_date" json:"completed_date"`
}

type IUserExamProcessUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	FetchByIdLesson(ctx context.Context, idLesson string) (Response, error)
	CreateOne(ctx context.Context, userProcess *UserExamProcess) error
	UpdateOne(ctx context.Context, userProcessID string, userProcess UserExamProcess) error
	UpsertOne(ctx context.Context, userProcessID string, userProcess *UserExamProcess) (Response, error)
	DeleteOne(ctx context.Context, userProcessID string) error
}
