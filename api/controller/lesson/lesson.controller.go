package lesson_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	lesson_domain "clean-architecture/domain/lesson"
)

type LessonController struct {
	LessonUseCase lesson_domain.ILessonUseCase
	AdminUseCase  admin_domain.IAdminUseCase
	Database      *bootstrap.Database
}
