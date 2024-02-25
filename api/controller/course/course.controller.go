package course_controller

import (
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
)

type CourseController struct {
	CourseUseCase course_domain.ICourseUseCase
	Database      *bootstrap.Database
}
