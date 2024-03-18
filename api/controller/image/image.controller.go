package image_controller

import (
	"clean-architecture/bootstrap"
	image_domain "clean-architecture/domain/image"
)

type ImageController struct {
	ImageUseCase image_domain.IImageUseCase
	Database     *bootstrap.Database
}
