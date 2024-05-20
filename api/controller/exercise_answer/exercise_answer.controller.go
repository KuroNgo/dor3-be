package exercise_answer_controller

import (
	"clean-architecture/bootstrap"
	exercise_answer_domain "clean-architecture/domain/exercise_answer"
	exercise_result_domain "clean-architecture/domain/exercise_result"
	user_domain "clean-architecture/domain/user"
	user_attempt_domain "clean-architecture/domain/user_attempt"
)

type ExerciseAnswerController struct {
	ExerciseAnswerUseCase exercise_answer_domain.IExerciseAnswerUseCase
	ExerciseResultUseCase exercise_result_domain.IExerciseResultUseCase
	UserAttemptUseCase    user_attempt_domain.IUserProcessUseCase
	UserUseCase           user_domain.IUserUseCase
	Database              *bootstrap.Database
}
