package quiz_result_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionQuizResult = "quiz_result"
)

type QuizResult struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuizID primitive.ObjectID `bson:"exam_id" json:"exam_id"`

	Score      int16     `bson:"score" json:"score"`
	StartedAt  time.Time `bson:"started_at" json:"started_at"`
	IsComplete int       `bson:"is_complete" json:"is_complete"`
}

type Response struct {
	Page       int64        `bson:"page" json:"page"`
	QuizResult []QuizResult `bson:"quiz_result" json:"quiz_result"`
	Statistics Statistics   `bson:"statistics" json:"statistics"`
}

type Statistics struct {
	TotalScore   int16   `bson:"total_score" json:"total_score"`
	AverageScore float64 `bson:"average_score" json:"average_score"`
	Percentage   float64 `bson:"percentage" json:"percentage"`
}

type IQuizResultRepository interface {
	FetchManyInUser(ctx context.Context, page string) (Response, error)
	FetchManyByQuizIDInUser(ctx context.Context, quizID string) (Response, error)
	GetResultsByUserIDAndQuizIDInUser(ctx context.Context, userID string, quizID string) (QuizResult, error)

	CreateOneInUser(ctx context.Context, quizResult *QuizResult) error
	DeleteOneInUser(ctx context.Context, quizResultID string) error
	UpdateStatusInUser(ctx context.Context, quizResultID string, status int) (*mongo.UpdateResult, error)

	CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int
	CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64
}
