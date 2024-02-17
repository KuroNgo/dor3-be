package quiz

import (
	quiz_domain "clean-architecture/domain/quiz"
	"time"
)

type quizUseCase struct {
	quizRepository quiz_domain.IQuizRepository
	contextTimeout time.Duration
}
