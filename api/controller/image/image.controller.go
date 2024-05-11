package image_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	image_domain "clean-architecture/domain/image"
	user_domain "clean-architecture/domain/user"
)

type ImageController struct {
	ImageUseCase image_domain.IImageUseCase
	AdminUseCase admin_domain.IAdminUseCase
	UserUseCase  user_domain.IUserUseCase
	Database     *bootstrap.Database
}
