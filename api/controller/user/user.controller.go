package user_controller

import (
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/user"
)

type UserController struct {
	UserUseCase user_domain.IUserUseCase
	Database    *bootstrap.Database
}
