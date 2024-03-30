package unit_controller

import (
	"clean-architecture/bootstrap"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
)

type UnitController struct {
	UnitUseCase unit_domain.IUnitUseCase
	UserUseCase user_domain.IUserUseCase
	Database    *bootstrap.Database
}
