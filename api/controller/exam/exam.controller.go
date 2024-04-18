package exam_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exam_domain "clean-architecture/domain/exam"
	user_domain "clean-architecture/domain/user"
)

type ExamsController struct {
	ExamUseCase  exam_domain.IExamUseCase
	UserUseCase  user_domain.IUserUseCase
	AdminUseCase admin_domain.IAdminUseCase
	Database     *bootstrap.Database
}
