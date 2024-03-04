package activity_log

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ActivityLog struct {
	LogID        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	ClientIP     string             `json:"client_ip" bson:"client_ip"`
	Method       string             `json:"method" bson:"method"`
	StatusCode   int                `json:"status_code" bson:"status_code"`
	BodySize     int                `json:"body_size" bson:"body_size"`
	Path         string             `json:"path" bson:"path"`
	Latency      string             `json:"latency" bson:"latency"`
	ActivityType string             `json:"activity_type" bson:"activity_type"`
	ActivityTime time.Time          `json:"activity_time" bson:"activity_time"`
}
