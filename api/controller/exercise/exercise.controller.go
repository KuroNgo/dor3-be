package exercise_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exercise_domain "clean-architecture/domain/exercise"
	user_domain "clean-architecture/domain/user"
)

type ExerciseController struct {
	ExerciseUseCase exercise_domain.IExerciseUseCase
	UserUseCase     user_domain.IUserUseCase
	AdminUseCase    admin_domain.IAdminUseCase
	Database        *bootstrap.Database
}
