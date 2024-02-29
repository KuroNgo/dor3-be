package internal

import (
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
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

func IsValidCourse(course course_domain.Input) error {
	if course.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if course.Description == "" {
		return fmt.Errorf("description cannot be empty")
	}

	if course.Level == 0 {
		return fmt.Errorf("level cannot be empty")
	}

	return nil
}

func IsValidLesson(lesson lesson_domain.Lesson) error {
	if lesson.Name == "" {
		return fmt.Errorf("name lesson cannot be empty")
	}
	if lesson.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}
