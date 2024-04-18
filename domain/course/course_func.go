package course_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
}

//go:generate mockery --name ICourseUseCase
type ICourseUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	UpdateOne(ctx context.Context, course Course) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, course *Course) error
	DeleteOne(ctx context.Context, courseID string) error
	CountCourse(ctx context.Context) int64
}
