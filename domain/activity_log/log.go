package activity_log_domain

import (
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
	Method       string             `json:"method" bson:"method"`
	StatusCode   int                `json:"status_code" bson:"status_code"`
	BodySize     int                `json:"body_size" bson:"body_size"`
	Path         string             `json:"path" bson:"path"`
	Latency      string             `json:"latency" bson:"latency"`
	Error        string             `json:"error" bson:"error"`
	ActivityTime time.Time          `json:"activity_time" bson:"activity_time"`
}

type Response struct {
	ActivityLog []ActivityLog
}

type IActivityRepository interface {
	CreateOne(ctx context.Context, log ActivityLog) error
	DeleteOne(ctx context.Context, logID string) error
	DeleteOneByTime(ctx context.Context, time time.Duration) error
	FetchMany(ctx context.Context, page string) (Response, error)
}

type IActivityRepositoryV2 interface {
	CreateOne(ctx context.Context, log ActivityLog) error
}
