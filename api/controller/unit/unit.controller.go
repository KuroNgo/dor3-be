package unit_controller

import (
	"clean-architecture/bootstrap"
	unit_domain "clean-architecture/domain/unit"
)

type UnitController struct {
	UnitUseCase unit_domain.IUnitUseCase
	Database    *bootstrap.Database
}
