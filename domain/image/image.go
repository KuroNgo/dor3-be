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
	Category  string             `bson:"category" json:"category"`
	AssetId   string             `bson:"asset_id" json:"asset_id"`
}

type Response struct {
	Image []Image
	Count int64 `json:"count"`
	Size  int64 `json:"size(KB)"`
}

//go:generate mockery --name IImageRepository
type IImageRepository interface {
	GetURLByName(ctx context.Context, name string) (Image, error)
	FetchMany(ctx context.Context) (Response, error)

	CreateOne(ctx context.Context, image *Image) error
	UpdateOne(ctx context.Context, imageID string, image Image) error

	DeleteOne(ctx context.Context, imageID string) error
}
