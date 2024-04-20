package exercise_quesiton_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	user_domain "clean-architecture/domain/user"
)

type ExerciseQuestionsController struct {
	ExerciseQuestionUseCase exercise_questions_domain.IExerciseQuestionUseCase
	UserUseCase             user_domain.IUserUseCase
	AdminUseCase            admin_domain.IAdminUseCase
	Database                *bootstrap.Database
}
