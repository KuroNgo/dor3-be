package exam_answer_controller

import (
	"clean-architecture/bootstrap"
	exam_answer_domain "clean-architecture/domain/exam_answer"
	exam_result_domain "clean-architecture/domain/exam_result"
	user_domain "clean-architecture/domain/user"
	user_attempt_domain "clean-architecture/domain/user_attempt"
)

type ExamAnswerController struct {
	ExamAnswerUseCase  exam_answer_domain.IExamAnswerUseCase
	ExamResultUseCase  exam_result_domain.IExamResultUseCase
	UserAttemptUseCase user_attempt_domain.IUserProcessUseCase
	UserUseCase        user_domain.IUserUseCase
	Database           *bootstrap.Database
}
