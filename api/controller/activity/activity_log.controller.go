package activity_controller

import (
	"clean-architecture/bootstrap"
	activity_log_domain "clean-architecture/domain/activity_log"
)

type ActivityController struct {
	ActivityUseCase activity_log_domain.IActivityUseCase
	Database        *bootstrap.Database
}

type ActivityControllerV2 struct {
	ActivityUseCase activity_log_domain.IActivityUseCaseV2
	Database        *bootstrap.Database
}
