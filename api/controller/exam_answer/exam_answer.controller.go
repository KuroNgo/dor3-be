package exam_answer

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exam_answer_domain "clean-architecture/domain/exam_answer"
	user_domain "clean-architecture/domain/user"
)

type ExamAnswerController struct {
	ExamAnswerUseCase exam_answer_domain.IExamAnswerUseCase
	AdminUseCase      admin_domain.IAdminUseCase
	UserUseCase       user_domain.IUserUseCase
	Database          *bootstrap.Database
}
