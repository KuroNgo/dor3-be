package image_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionImage = "image"
)

type Image struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	ImageName string             `bson:"image_name" json:"image_name"`
	Size      int64              `bson:"size" json:"size"`
	ImageUrl  string             `bson:"image_url" json:"image_url"`
}

type Response struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	ImageName string             `bson:"image_name" json:"image_name"`
	Size      int64              `bson:"size" json:"size"`
	ImageUrl  string             `bson:"image_url" json:"image_url"`
}

//go:generate mockery --name IAudioRepository
type IImageRepository interface {
	GetURLByName(ctx context.Context, name string) (Image, error)
	FetchMany(ctx context.Context) ([]Image, error)
	UpdateOne(ctx context.Context, imageID string, image Image) error
	CreateOne(ctx context.Context, image *Image) error
	DeleteOne(ctx context.Context, imageID string) error
	DeleteMany(ctx context.Context, imageID ...string) error
}
