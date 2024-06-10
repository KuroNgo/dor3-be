package exam_result_controller

import (
	"clean-architecture/bootstrap"
	exam_result_domain "clean-architecture/domain/exam_result"
	user_domain "clean-architecture/domain/user"
	user_attempt_domain "clean-architecture/domain/user_process"
)

type ExamResultController struct {
	ExamResultUseCase  exam_result_domain.IExamResultUseCase
	UserAttemptUseCase user_attempt_domain.IUserProcessUseCase
	UserUseCase        user_domain.IUserUseCase
	Database           *bootstrap.Database
}
