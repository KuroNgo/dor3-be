package exam_answer_controller

import (
	"clean-architecture/bootstrap"
	exam_answer_domain "clean-architecture/domain/exam_answer"
	user_domain "clean-architecture/domain/user"
)

type ExamAnswerController struct {
	ExamAnswerUseCase exam_answer_domain.IExamAnswerUseCase
	UserUseCase       user_domain.IUserUseCase
	Database          *bootstrap.Database
}
