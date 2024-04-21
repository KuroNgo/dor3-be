package quiz_answer_controller

import (
	"clean-architecture/bootstrap"
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	user_domain "clean-architecture/domain/user"
)

type QuizAnswerController struct {
	QuizAnswerUseCase quiz_answer_domain.IQuizAnswerUseCase
	UserUseCase       user_domain.IUserUseCase
	Database          *bootstrap.Database
}
