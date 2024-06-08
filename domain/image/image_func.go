package image_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	ImageName string             `bson:"image_name" json:"image_name"`
}

//go:generate mockery --name IImageUseCase
type IImageUseCase interface {
	GetURLByName(ctx context.Context, name string) (Image, error)
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchByCategory(ctx context.Context, category string, page string) (Response, error)

	CreateOne(ctx context.Context, image *Image) error
	UpdateOne(ctx context.Context, image *Image) error
	DeleteOne(ctx context.Context, imageID string) error
}
