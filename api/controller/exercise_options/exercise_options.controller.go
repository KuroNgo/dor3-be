package exercise_options_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exercise_options_domain "clean-architecture/domain/exercise_options"
	user_domain "clean-architecture/domain/user"
)

type ExerciseOptionsController struct {
	ExerciseOptionsUseCase exercise_options_domain.IExerciseOptionUseCase
	UserUseCase            user_domain.IUserUseCase
	AdminUseCase           admin_domain.IAdminUseCase
	Database               *bootstrap.Database
}
