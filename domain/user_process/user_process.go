package user_process

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUserProcess = "user_process"
)

type UserProcess struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	CourseID      primitive.ObjectID `bson:"course_id" json:"course_id"`
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	Score         int                `bson:"score" json:"score"`
	ProcessStatus int                `bson:"process_status" json:"process_status"`
	CompletedDate time.Time          `bson:"completed_date" json:"completed_date"`
}

type Response struct {
	UserProcess                []UserProcess `bson:"user_process" json:"user_process"`
	StudyTime                  time.Time     `bson:"study_time" json:"study_time"`
	TotalNumberOfLessonLearned int16         `bson:"total_number_of_lesson_learned" json:"total_number_of_lesson_learned"`
	TotalScore                 int64         `bson:"total_score" json:"total_score"`
	AverageScore               int8          `bson:"average_score" json:"average_score"`
}

type IUserProcessRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	FetchByIdLesson(ctx context.Context, idLesson string) (Response, error)
	CreateOne(ctx context.Context, userProcess *UserProcess) error
	UpdateOne(ctx context.Context, userProcessID string, userProcess UserProcess) error
	UpsertOne(ctx context.Context, userProcessID string, userProcess *UserProcess) (Response, error)
	DeleteOne(ctx context.Context, userProcessID string) error
}
