package lesson_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	course_domain "clean-architecture/domain/course"
	image_domain "clean-architecture/domain/image"
	lesson_domain "clean-architecture/domain/lesson"
	user_domain "clean-architecture/domain/user"
)

type LessonController struct {
	LessonUseCase lesson_domain.ILessonUseCase
	CourseUseCase course_domain.ICourseUseCase
	ImageUseCase  image_domain.IImageUseCase
	AdminUseCase  admin_domain.IAdminUseCase
	UserUseCase   user_domain.IUserUseCase
	Database      *bootstrap.Database
}
