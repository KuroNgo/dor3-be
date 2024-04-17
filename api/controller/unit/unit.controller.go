package unit_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	unit_domain "clean-architecture/domain/unit"
)

type UnitController struct {
	UnitUseCase  unit_domain.IUnitUseCase
	AdminUseCase admin_domain.IAdminUseCase
	Database     *bootstrap.Database
}
