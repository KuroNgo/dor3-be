package unit_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
)

type UnitController struct {
	UnitUseCase  unit_domain.IUnitUseCase
	UserUseCase  user_domain.IUserUseCase
	AdminUseCase admin_domain.IAdminUseCase
	Database     *bootstrap.Database
}
