package lesson_controller

import (
	"clean-architecture/bootstrap"
	lesson_domain "clean-architecture/domain/lesson"
)

type LessonController struct {
	LessonUseCase lesson_domain.ILessonUseCase
	Database      *bootstrap.Database
}
