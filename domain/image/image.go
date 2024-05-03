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
	Page  int64   `json:"page"`
	Image []Image `json:"image" bson:"image"`
}

type Statistics struct {
	Count           int64 `json:"count"`
	SizeKB          int64 `json:"size(KB)"`
	SizeRemainingKB int64 `json:"size_remaining(KB)"`
	MaxSizeKB       int64 `json:"max_size(KB)"`
	SizeMB          int64 `json:"size(MB)"`
	SizeRemainingMB int64 `json:"size_remaining(MB)"`
	MaxSizeMB       int64 `json:"max_size(MB)"`
}

//go:generate mockery --name IImageRepository
type IImageRepository interface {
	GetURLByName(ctx context.Context, name string) (Image, error)
	FetchMany(ctx context.Context, page string) (Response, error)

	CreateOne(ctx context.Context, image *Image) error
	UpdateOne(ctx context.Context, imageID string, image *Image) error

	DeleteOne(ctx context.Context, imageID string) error
}
