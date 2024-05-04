package image_usecase

import (
	image_domain "clean-architecture/domain/image"
	"context"
	"time"
)

type imageUseCase struct {
	imageRepository image_domain.IImageRepository
	contextTimeout  time.Duration
}

func (i *imageUseCase) FetchByCategory(ctx context.Context, category string, page string) (image_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, i.contextTimeout)
	defer cancel()

	image, err := i.imageRepository.FetchByCategory(ctx, category, page)
	if err != nil {
		return image_domain.Response{}, err
	}

	return image, err
}

func NewImageUseCase(imageRepository image_domain.IImageRepository, timeout time.Duration) image_domain.IImageUseCase {
	return &imageUseCase{
		imageRepository: imageRepository,
		contextTimeout:  timeout,
	}
}
func (i *imageUseCase) GetURLByName(ctx context.Context, name string) (image_domain.Image, error) {
	ctx, cancel := context.WithTimeout(ctx, i.contextTimeout)
	defer cancel()

	image, err := i.imageRepository.GetURLByName(ctx, name)
	if err != nil {
		return image_domain.Image{}, err
	}

	return image, err
}

func (i *imageUseCase) FetchMany(ctx context.Context, page string) (image_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, i.contextTimeout)
	defer cancel()

	quiz, err := i.imageRepository.FetchMany(ctx, page)
	if err != nil {
		return image_domain.Response{}, err
	}

	return quiz, err
}

func (i *imageUseCase) UpdateOne(ctx context.Context, imageID string, image *image_domain.Image) error {
	ctx, cancel := context.WithTimeout(ctx, i.contextTimeout)
	defer cancel()

	err := i.imageRepository.UpdateOne(ctx, imageID, image)
	if err != nil {
		return err
	}

	return nil
}

func (i *imageUseCase) CreateOne(ctx context.Context, image *image_domain.Image) error {
	ctx, cancel := context.WithTimeout(ctx, i.contextTimeout)
	defer cancel()
	err := i.imageRepository.CreateOne(ctx, image)

	if err != nil {
		return err
	}

	return nil
}

func (i *imageUseCase) DeleteOne(ctx context.Context, imageID string) error {
	ctx, cancel := context.WithTimeout(ctx, i.contextTimeout)
	defer cancel()

	err := i.imageRepository.DeleteOne(ctx, imageID)
	if err != nil {
		return err
	}

	return err
}
