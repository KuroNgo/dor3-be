package course_controller

import (
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	mean_domain "clean-architecture/domain/mean"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
)

type CourseController struct {
	CourseUseCase     course_domain.ICourseUseCase
	LessonUseCase     lesson_domain.ILessonUseCase
	UnitUseCase       unit_domain.IUnitUseCase
	VocabularyUseCase vocabulary_domain.IVocabularyUseCase
	MeanUseCase       mean_domain.IMeanUseCase
	Database          *bootstrap.Database
}
