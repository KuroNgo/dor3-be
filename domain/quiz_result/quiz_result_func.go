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
	FetchManyInUser(ctx context.Context, page string) (Response, error)
	FetchManyByQuizIDInUser(ctx context.Context, quizID string) (Response, error)
	GetResultsByUserIDAndQuizIDInUser(ctx context.Context, userID string, quizID string) (QuizResult, error)

	CreateOneInUser(ctx context.Context, quizResult *QuizResult) error
	DeleteOneInUser(ctx context.Context, quizResultID string) error
}
