package quiz_question_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	quiz_question_domain "clean-architecture/domain/quiz_question"
	user_domain "clean-architecture/domain/user"
)

type QuizQuestionsController struct {
	QuizQuestionUseCase quiz_question_domain.IQuizQuestionUseCase
	UserUseCase         user_domain.IUserUseCase
	AdminUseCase        admin_domain.IAdminUseCase
	Database            *bootstrap.Database
}
