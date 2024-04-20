package exam_question_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionExamQuestion = "exam_question"
)

type ExamQuestion struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	ExamID       primitive.ObjectID `bson:"exam_id" json:"exam_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`

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
	Count        int `bson:"count" json:"count"`
	ExamQuestion []ExamQuestion
}

type IExamQuestionRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string) (Response, error)

	CreateOne(ctx context.Context, examQuestion *ExamQuestion) error
	UpdateOne(ctx context.Context, examQuestion *ExamQuestion) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, examID string) error
}
