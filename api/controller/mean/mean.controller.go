package mean_controller

import (
	"clean-architecture/bootstrap"
	mean_domain "clean-architecture/domain/mean"
)

type MeanController struct {
	MeanUseCase mean_domain.IMeanUseCase
	Database    *bootstrap.Database
}
