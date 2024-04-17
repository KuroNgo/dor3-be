package course_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionCourse = "course"
)

type Course struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdated  string             `bson:"who_updated" json:"who_updated"`
}

type Response struct {
	Course []Course
}

//go:generate mockery --name ICourseRepository
type ICourseRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	UpdateOne(ctx context.Context, course Course) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, course *Course) error
	DeleteOne(ctx context.Context, courseID string) error
}
