package exam_result_controller

import (
	"clean-architecture/bootstrap"
	exam_result_domain "clean-architecture/domain/exam_result"
	user_domain "clean-architecture/domain/user"
)

type ExamResultController struct {
	ExamResultUseCase exam_result_domain.IExamResultUseCase
	UserUseCase       user_domain.IUserUseCase
	Database          *bootstrap.Database
}
