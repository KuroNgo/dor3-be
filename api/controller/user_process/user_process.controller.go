package user_attempt_controller

import (
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/user"
	user_attempt_domain "clean-architecture/domain/user_process"
)

type UserAttemptController struct {
	UserAttemptUseCase user_attempt_domain.IUserProcessUseCase
	UserUseCase        user_domain.IUserUseCase
	Database           *bootstrap.Database
}
