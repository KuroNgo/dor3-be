package course_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
)

type CourseController struct {
	CourseUseCase     course_domain.ICourseUseCase
	LessonUseCase     lesson_domain.ILessonUseCase
	UnitUseCase       unit_domain.IUnitUseCase
	VocabularyUseCase vocabulary_domain.IVocabularyUseCase
	UserUseCase       user_domain.IUserUseCase
	AdminUseCase      admin_domain.IAdminUseCase
	Database          *bootstrap.Database
}
