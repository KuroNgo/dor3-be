package user_exercise_attempt

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUserExerciseAttempt = "user_exercise_process"
)

type UserExerciseProcess struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	QuizID        primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	Score         int                `bson:"score" json:"score"`
	ProcessStatus int                `bson:"process_status" json:"process_status"`
	CompletedDate time.Time          `bson:"completed_date" json:"completed_date"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type Response struct {
	UserExerciseProcess []UserExerciseProcess
}

type IUserProcessRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	FetchByIdLesson(ctx context.Context, idLesson string) (Response, error)
	CreateOne(ctx context.Context, userProcess *UserExerciseProcess) error
	UpdateOne(ctx context.Context, userProcessID string, userProcess UserExerciseProcess) error
	UpsertOne(ctx context.Context, userProcessID string, userProcess *UserExerciseProcess) (Response, error)
	DeleteOne(ctx context.Context, userProcessID string) error
}
