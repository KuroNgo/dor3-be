package lesson_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	lesson_domain "clean-architecture/domain/lesson"
	user_domain "clean-architecture/domain/user"
)

type LessonController struct {
	LessonUseCase lesson_domain.ILessonUseCase
	UserUseCase   user_domain.IUserUseCase
	AdminUseCase  admin_domain.IAdminUseCase
	Database      *bootstrap.Database
}
