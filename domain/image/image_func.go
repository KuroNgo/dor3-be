package image_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AutoMatch struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	ImageName string             `bson:"image_name" json:"image_name"`
	Size      int64              `bson:"size" json:"size"`
	ImageUri  string             `bson:"image_uri" json:"image_uri"`
}

type Input struct {
	ImageName string `bson:"image_name" json:"image_name"`
}

type IImageUseCase interface {
	GetURLByName(ctx context.Context, name string) (Image, error)
	FetchMany(ctx context.Context) ([]Image, error)
	UpdateOne(ctx context.Context, imageID string, image Image) error
	CreateOne(ctx context.Context, image *Image) error
	DeleteOne(ctx context.Context, imageID string) error
	DeleteMany(ctx context.Context, imageID ...string) error
}
