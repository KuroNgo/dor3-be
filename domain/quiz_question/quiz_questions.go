package quiz_question

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionQuizQuestion = "quiz_question"
)

type QuizQuestion struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	QuizID primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"`
	Level   int    `bson:"level" json:"level"`

	// admin add metadata of file and system will be found it
	Filename      string `bson:"filename" json:"filename"`
	AudioDuration string `bson:"audio_duration" json:"audio_duration"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	QuizQuestion []QuizQuestion
}

type IQuizQuestionRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, quizID string) (Response, error)
	UpdateOne(ctx context.Context, quizQuestionID string, quizQuestion QuizQuestion) error
	CreateOne(ctx context.Context, quizQuestion *QuizQuestion) error
	DeleteOne(ctx context.Context, quizID string) error
}
