package exam

import (
	"clean-architecture/bootstrap"
	exercise_domain "clean-architecture/domain/exercise"
	user_domain "clean-architecture/domain/user"
)

type ExamController struct {
	ExamUseCase exercise_domain.IExerciseUseCase
	UserUseCase user_domain.IUserUseCase
	Database    *bootstrap.Database
}
