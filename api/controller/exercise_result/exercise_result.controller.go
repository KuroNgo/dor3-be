package exercise_result_controller

import (
	"clean-architecture/bootstrap"
	exercise_result_domain "clean-architecture/domain/exercise_result"
	user_domain "clean-architecture/domain/user"
)

type ExerciseResultController struct {
	ExerciseResultUseCase exercise_result_domain.IExerciseResultUseCase
	UserUseCase           user_domain.IUserUseCase
	Database              *bootstrap.Database
}
