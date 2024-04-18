package quiz_result_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Auto struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuizID primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Score     int16     `bson:"score" json:"score"`
	StartedAt time.Time `bson:"started_at" json:"started_at"`
	Status    int       `bson:"status" json:"status"`
}

type IQuizResultUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByQuizID(ctx context.Context, quizID string) (Response, error)

	CreateOne(ctx context.Context, quizResult *QuizResult) error
	DeleteOne(ctx context.Context, quizResultID string) error
}
