package unit_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
)

type UnitController struct {
	UnitUseCase   unit_domain.IUnitUseCase
	LessonUseCase lesson_domain.ILessonUseCase
	UserUseCase   user_domain.IUserUseCase
	AdminUseCase  admin_domain.IAdminUseCase
	Database      *bootstrap.Database
}
