package exercise_questions

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionExerciseQuestion = "exercise_question"
)

type ExerciseQuestion struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	ExerciseID primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`

	// admin add metadata of file and system will be found it
	Filename      string `bson:"filename" json:"filename"`
	AudioDuration string `bson:"audio_duration" json:"audio_duration"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	ExerciseQuestion []ExerciseQuestion
}

type IExerciseQuestionRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, exerciseID string) (Response, error)
	UpdateOne(ctx context.Context, exerciseQuestionID string, exerciseQuestion ExerciseQuestion) error
	CreateOne(ctx context.Context, exerciseQuestion *ExerciseQuestion) error
	DeleteOne(ctx context.Context, exerciseID string) error
}
