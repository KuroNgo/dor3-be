package course_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Course []Course `bson:"data" json:"data"`
}

//go:generate mockery --name ICourseRepository
type ICourseRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	UpdateOne(ctx context.Context, courseID string, course Course) error
	CreateOne(ctx context.Context, course *Course) error
	UpsertOne(ctx context.Context, id string, course *Course) (*Response, error)
	DeleteOne(ctx context.Context, courseID string) error
}
