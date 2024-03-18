package user_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
)

type UserController struct {
	UserUseCase user_domain.IUserUseCase
	Database    *bootstrap.Database
}

type LoginFromRoleController struct {
	UserUseCase  user_domain.IUserUseCase
	AdminUseCase admin_domain.IAdminUseCase
	Database     *bootstrap.Database
}
