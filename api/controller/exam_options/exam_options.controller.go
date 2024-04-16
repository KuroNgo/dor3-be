package exam_options_controller

import (
	"clean-architecture/bootstrap"
	exam_options_domain "clean-architecture/domain/exam_options"
	user_domain "clean-architecture/domain/user"
)

type ExamOptionsController struct {
	ExamOptionsUseCase exam_options_domain.IExamOptionsUseCase
	UserUseCase        user_domain.IUserUseCase
	Database           *bootstrap.Database
}
