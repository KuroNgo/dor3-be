package feedback_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Input struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title         string             `bson:"title" json:"title"`
	Content       string             `bson:"content" json:"content"`
	SubmittedDate time.Time          `bson:"submitted_date" json:"submitted_date"`
}

type IFeedbackUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	FetchByUserID(ctx context.Context, userID string) (Response, error)
	FetchBySubmittedDate(ctx context.Context, userID string) (Response, error)
	CreateOneByUser(ctx context.Context, feedback *Feedback) error
	DeleteOneByAdmin(ctx context.Context, feedbackID string) error
}
