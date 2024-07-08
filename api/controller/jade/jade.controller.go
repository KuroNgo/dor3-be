package jade_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	jade_domain "clean-architecture/domain/jade"
	user_domain "clean-architecture/domain/user"
)

type JadeController struct {
	JadeUseCase  jade_domain.IJadeUseCase
	UserUseCase  user_domain.IUserUseCase
	AdminUseCase admin_domain.IAdminUseCase
	Database     *bootstrap.Database
}
