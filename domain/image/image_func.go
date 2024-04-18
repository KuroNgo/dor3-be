package image_domain

import (
	"context"
)

type Input struct {
	ImageName string `bson:"image_name" json:"image_name"`
}

//go:generate mockery --name IImageUseCase
type IImageUseCase interface {
	GetURLByName(ctx context.Context, name string) (Image, error)
	FetchMany(ctx context.Context) (Response, error)

	CreateOne(ctx context.Context, image *Image) error
	UpdateOne(ctx context.Context, imageID string, image *Image) error
	DeleteOne(ctx context.Context, imageID string) error
}
