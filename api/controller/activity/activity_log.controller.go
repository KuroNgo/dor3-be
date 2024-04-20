package activity_controller

import (
	"clean-architecture/bootstrap"
	activity_log_domain "clean-architecture/domain/activity_log"
	admin_domain "clean-architecture/domain/admin"
)

type ActivityController struct {
	ActivityUseCase activity_log_domain.IActivityUseCase
	AdminUseCase    admin_domain.IAdminUseCase
	Database        *bootstrap.Database
}
