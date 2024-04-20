package exercise_answer_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exercise_answer_domain "clean-architecture/domain/exercise_answer"
	user_domain "clean-architecture/domain/user"
)

type ExerciseAnswerController struct {
	ExerciseAnswerUseCase exercise_answer_domain.IExerciseAnswerUseCase
	AdminUseCase          admin_domain.IAdminUseCase
	UserUseCase           user_domain.IUserUseCase
	Database              *bootstrap.Database
}
