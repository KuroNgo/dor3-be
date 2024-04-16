package quiz_answer

import (
	exercise_result_domain "clean-architecture/domain/exercise_result"
	"time"
)

type exerciseResultUseCase struct {
	exerciseQuestionRepository exercise_result_domain.IExerciseResultRepository
	contextTimeout             time.Duration
}
