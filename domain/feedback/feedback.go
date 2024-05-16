package feedback_domain

import (
	user_domain "clean-architecture/domain/user"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionFeedback = "feedback"
)

type Feedback struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title         string             `bson:"title" json:"title"`
	Content       string             `bson:"content" json:"content"`
	Feeling       string             `bson:"feeling" json:"feeling"`
	SubmittedDate time.Time          `bson:"submitted_date" json:"submitted_date"`
	IsLoveWeb     int                `bson:"is_love_web" json:"is_love_web"`
	IsSeen        int                `bson:"is_seen" json:"is_seen"`
	SeenAt        time.Time          `bson:"seen_at" json:"seen_at"`
}

type FeedbackResponse struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	User          user_domain.User   `bson:"user" json:"user"`
	Title         string             `bson:"title" json:"title"`
	Content       string             `bson:"content" json:"content"`
	Feeling       string             `bson:"feeling" json:"feeling"`
	SubmittedDate time.Time          `bson:"submitted_date" json:"submitted_date"`
	IsLoveWeb     int                `bson:"is_love_web" json:"is_love_web"`
	IsSeen        int                `bson:"is_seen" json:"is_seen"`
	SeenAt        time.Time          `bson:"seen_at" json:"seen_at"`
}

type Response struct {
	Page        int64              `bson:"page" json:"page"`
	CurrentPage int64              `bson:"current_page" json:"current_page"`
	Statistics  Statistics         `bson:"statistics" json:"statistics"`
	Feedback    []FeedbackResponse `json:"feedback" bson:"feedback"`
}

type Statistics struct {
	Total             int64   `bson:"total" json:"total"`
	TotalFeeling      int32   `bson:"feeling" json:"feeling"`
	TotalIsSeen       int32   `bson:"is_seen" json:"is_seen"`
	TotalIsNotSeen    int32   `bson:"is_not_seen" json:"is_not_seen"`
	TotalIsLoveWeb    int32   `bson:"is_love_web" json:"is_love_web"`
	CountSad          float32 `bson:"count_sad" json:"count_sad"`
	CountHappy        float32 `bson:"count_happy" json:"count_happy"`
	CountDisappointed float32 `bson:"count_disappointed" json:"count_disappointed"`
	CountGood         float32 `bson:"count_good" json:"count_good"`
}

type IFeedbackRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchByUserID(ctx context.Context, userID string, page string) (Response, error)
	FetchBySubmittedDate(ctx context.Context, date string, page string) (Response, error)

	CreateOneByUser(ctx context.Context, feedback *Feedback) error
	DeleteOneByAdmin(ctx context.Context, feedbackID string) error
	UpdateSeen(ctx context.Context, id string, isSeen int) error
	Statistics(ctx context.Context) (Statistics, error)
}
