package exam_controller

import (
	"clean-architecture/bootstrap"
	exam_domain "clean-architecture/domain/exam"
	user_domain "clean-architecture/domain/user"
)

type ExamsController struct {
	ExamUseCase exam_domain.IExamUseCase
	UserUseCase user_domain.IUserUseCase
	Database    *bootstrap.Database
}
