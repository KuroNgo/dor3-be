package activity_log_domain

import (
	admin_domain "clean-architecture/domain/admin"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionActivityLog = "activity_log"
)

// ActivityLog this code process write log user automation
type ActivityLog struct {
	LogID        primitive.ObjectID `json:"_id" bson:"_id"`
	ClientIP     string             `json:"client_ip" bson:"client_ip"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	Method       string             `json:"method" bson:"method"`
	StatusCode   int                `json:"status_code" bson:"status_code"`
	BodySize     int                `json:"body_size" bson:"body_size"`
	Path         string             `json:"path" bson:"path"`
	Latency      string             `json:"latency" bson:"latency"`
	Error        string             `json:"error" bson:"error"`
	ActivityTime time.Time          `json:"activity_time" bson:"activity_time"`
	ExpireAt     time.Time          `json:"expire_at" bson:"expire_at"`
}

// ActivityLog this code process write log user automation
type ActivityLogResponse struct {
	LogID        primitive.ObjectID `json:"_id" bson:"_id"`
	ClientIP     string             `json:"client_ip" bson:"client_ip"`
	UserID       admin_domain.Admin `json:"user_id" bson:"user_id"`
	Method       string             `json:"method" bson:"method"`
	StatusCode   int                `json:"status_code" bson:"status_code"`
	BodySize     int                `json:"body_size" bson:"body_size"`
	Path         string             `json:"path" bson:"path"`
	Latency      string             `json:"latency" bson:"latency"`
	Error        string             `json:"error" bson:"error"`
	ActivityTime time.Time          `json:"activity_time" bson:"activity_time"`
	ExpireAt     time.Time          `json:"expire_at" bson:"expire_at"`
}

type Response struct {
	PageCurrent int64         `json:"page_current" bson:"page_current"`
	Page        int64         `json:"page" bson:"page"`
	ActivityLog []ActivityLog `json:"activity_log" bson:"activity_log"`
	Statistics  Statistics    `json:"statistics" bson:"statistics"`
}

type Statistics struct {
	Total int64 `json:"total" bson:"total"`
}

type IActivityRepository interface {
	CreateOne(ctx context.Context, log ActivityLog) error
	FetchMany(ctx context.Context, page string) (Response, error)
	Statistics(ctx context.Context) (Statistics, error)
}
