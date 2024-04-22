package quiz_options_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	quiz_options_domain "clean-architecture/domain/quiz_options"
	user_domain "clean-architecture/domain/user"
)

type QuizOptionsController struct {
	QuizOptionsUseCase quiz_options_domain.IQuizOptionUseCase
	UserUseCase        user_domain.IUserUseCase
	AdminUseCase       admin_domain.IAdminUseCase
	Database           *bootstrap.Database
}
