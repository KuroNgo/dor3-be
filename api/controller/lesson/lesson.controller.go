package lesson_controller

import (
	"clean-architecture/bootstrap"
	lesson_domain "clean-architecture/domain/lesson"
	user_domain "clean-architecture/domain/user"
)

type LessonController struct {
	LessonUseCase lesson_domain.ILessonUseCase
	UserUseCase   user_domain.IUserUseCase
	Database      *bootstrap.Database
}
