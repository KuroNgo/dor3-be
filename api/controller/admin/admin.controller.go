package admin_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
)

type AdminController struct {
	AdminUseCase admin_domain.IAdminUseCase
	Database     *bootstrap.Database
}
