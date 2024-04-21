package quiz_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	quiz_domain "clean-architecture/domain/quiz"
	user_domain "clean-architecture/domain/user"
)

type QuizController struct {
	QuizUseCase  quiz_domain.IQuizUseCase
	UserUseCase  user_domain.IUserUseCase
	AdminUseCase admin_domain.IAdminUseCase
	Database     *bootstrap.Database
}
