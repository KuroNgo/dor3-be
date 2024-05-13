package activity_controller

import (
	"clean-architecture/bootstrap"
	activity_log_domain "clean-architecture/domain/activity_log"
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
)

type ActivityController struct {
	ActivityUseCase activity_log_domain.IActivityUseCase
	AdminUseCase    admin_domain.IAdminUseCase
	UserUseCase     user_domain.IUserUseCase
	Database        *bootstrap.Database
}
