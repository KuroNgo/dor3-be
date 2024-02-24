package internal

import (
	quiz_domain "clean-architecture/domain/quiz"
	"fmt"
)

func IsValidQuiz(quiz quiz_domain.Input) error {
	if quiz.CorrectAnswer == "" {
		return fmt.Errorf("correct answer cannot be empty")
	}

	if quiz.Explanation == "" {
		return fmt.Errorf("explanation cannot be empty")
	}
	return nil
}
