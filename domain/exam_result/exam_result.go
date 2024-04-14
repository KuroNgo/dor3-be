package exam_result_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionExamResult = "exam_result"
)

type ExamResult struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExamID primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Score     int16     `bson:"score" json:"score"`
	StartedAt time.Time `bson:"started_at" json:"started_at"`
	Status    int       `bson:"status" json:"status"`
}

type Response struct {
	ExamResult []ExamResult
	TotalScore int16 `bson:"total_score" json:"total_score"`
}

type IExamOptionsRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByExamID(ctx context.Context, examID string) (Response, error)
	CreateOne(ctx context.Context, examResult *ExamResult) error
	DeleteOne(ctx context.Context, examResultID string) error
}
