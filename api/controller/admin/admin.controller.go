package admin_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
)

type AdminController struct {
	AdminUseCase admin_domain.IAdminUseCase
	UserUseCase  user_domain.IUserUseCase
	Database     *bootstrap.Database
}
