package quiz_answer_controller

import (
	"clean-architecture/bootstrap"
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	quiz_result_domain "clean-architecture/domain/quiz_result"
	user_domain "clean-architecture/domain/user"
	user_attempt_domain "clean-architecture/domain/user_process"
)

type QuizAnswerController struct {
	QuizAnswerUseCase  quiz_answer_domain.IQuizAnswerUseCase
	QuizResultUseCase  quiz_result_domain.IQuizResultUseCase
	UserAttemptUseCase user_attempt_domain.IUserProcessUseCase
	UserUseCase        user_domain.IUserUseCase
	Database           *bootstrap.Database
}
