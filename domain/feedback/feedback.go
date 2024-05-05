package feedback_domain

import (
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
}

type Response struct {
	Page     int64      `bson:"page" json:"page"`
	Total    int64      `bson:"total" json:"total"`
	Feedback []Feedback `json:"feedback" bson:"feedback"`
}

type IFeedbackRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchByUserID(ctx context.Context, userID string, page string) (Response, error)
	FetchBySubmittedDate(ctx context.Context, date string, page string) (Response, error)

	CreateOneByUser(ctx context.Context, feedback *Feedback) error
	DeleteOneByAdmin(ctx context.Context, feedbackID string) error
}

//API Feedback trong có field cảm xúc (thất vọng, tạm được, hài lòng, quá tuyệt)
//hoặc lưu theo cách nào tối ưu cũng được.
//Field ý kiến và field (có muố ở lại trang web không) trả về 0,1 hay true false cũng dc
