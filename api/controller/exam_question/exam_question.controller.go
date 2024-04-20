package exam_question_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exam_question_domain "clean-architecture/domain/exam_question"
	user_domain "clean-architecture/domain/user"
)

type ExamQuestionsController struct {
	ExamQuestionUseCase exam_question_domain.IExamQuestionUseCase
	UserUseCase         user_domain.IUserUseCase
	AdminUseCase        admin_domain.IAdminUseCase
	Database            *bootstrap.Database
}
