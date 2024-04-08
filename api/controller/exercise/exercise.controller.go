package exercise_controller

import (
	"clean-architecture/bootstrap"
	exercise_domain "clean-architecture/domain/exercise"
	user_domain "clean-architecture/domain/user"
)

type ExerciseController struct {
	ExerciseUseCase exercise_domain.IExerciseUseCase
	UserUseCase     user_domain.IUserUseCase
	Database        *bootstrap.Database
}
