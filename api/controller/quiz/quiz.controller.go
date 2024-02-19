package quiz_controller

import (
	"clean-architecture/bootstrap"
	quiz_domain "clean-architecture/domain/quiz"
)

type QuizController struct {
	QuizUseCase quiz_domain.IQuizUseCase
	Database    *bootstrap.Database
}
