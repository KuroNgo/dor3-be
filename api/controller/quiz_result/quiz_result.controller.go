package quiz_result_controller

import (
	"clean-architecture/bootstrap"
	quiz_result_domain "clean-architecture/domain/quiz_result"
	user_domain "clean-architecture/domain/user"
)

type QuizResultController struct {
	QuizResultUseCase quiz_result_domain.IQuizResultUseCase
	UserUseCase       user_domain.IUserUseCase
	Database          *bootstrap.Database
}
