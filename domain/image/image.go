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
	ImageName string             `bson:"image-name" json:"image-name"`
	Size      int64              `bson:"size" json:"size"`
	ImageUri  string             `bson:"image-uri" json:"image-uri"`
}

type Response struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	ImageName string             `bson:"image-name" json:"image-name"`
	Size      int64              `bson:"size" json:"size"`
	ImageUri  string             `bson:"image-uri" json:"image-uri"`
}

//go:generate mockery --name IAudioRepository
type IImageRepository interface {
	FetchMany(ctx context.Context) ([]Image, error)
	UpdateOne(ctx context.Context, imageID string, image Image) error
	CreateOne(ctx context.Context, image *Image) error
	DeleteOne(ctx context.Context, imageID string) error
	DeleteMany(ctx context.Context, imageID ...string) error
}
