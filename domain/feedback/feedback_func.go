package feedback_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Feeling   string             `bson:"feeling" json:"feeling"`
	IsLoveWeb int                `bson:"is_love_web" json:"is_love_web"`
	IsSeen    int                `bson:"is_seen" json:"is_seen"`
}

type IFeedbackUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchByUserID(ctx context.Context, userID string, page string) (Response, error)
	FetchBySubmittedDate(ctx context.Context, date string, page string) (Response, error)
	CreateOneByUser(ctx context.Context, feedback *Feedback) error
	DeleteOneByAdmin(ctx context.Context, feedbackID string) error
}
