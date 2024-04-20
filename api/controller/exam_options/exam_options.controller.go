package exam_options_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exam_options_domain "clean-architecture/domain/exam_options"
	user_domain "clean-architecture/domain/user"
)

type ExamOptionsController struct {
	ExamOptionsUseCase exam_options_domain.IExamOptionsUseCase
	UserUseCase        user_domain.IUserUseCase
	AdminUseCase       admin_domain.IAdminUseCase
	Database           *bootstrap.Database
}
