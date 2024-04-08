package quiz_controller

import (
	"clean-architecture/bootstrap"
	quiz_domain "clean-architecture/domain/quiz"
	user_domain "clean-architecture/domain/user"
)

type QuizController struct {
	QuizUseCase quiz_domain.IQuizUseCase
	UserUseCase user_domain.IUserUseCase
	Database    *bootstrap.Database
}
