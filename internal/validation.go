package internal

import (
	unit_domain "clean-architecture/domain/_unit"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	quiz_domain "clean-architecture/domain/quiz"
	"gopkg.in/errgo.v2/fmt/errors"
)

func IsValidQuiz(quiz quiz_domain.Input) error {
	if quiz.Question == "" {
		return errors.New("question cannot be empty")
	}
	if len(quiz.Options) == 0 {
		return errors.New("options cannot be empty")
	}

	if quiz.CorrectAnswer == "" {
		return errors.New("correct answer cannot be empty")
	}

	if quiz.Explanation == "" {
		return errors.New("explanation cannot be empty")
	}

	if quiz.QuestionType == "" {
		return errors.New("question type cannot be empty")

	}
	return nil
}

func IsValidCourse(course course_domain.Input) error {
	if course.Name == "" {
		return errors.New("name cannot be empty")
	}

	if course.Description == "" {
		return errors.New("description cannot be empty")
	}

	return nil
}

func IsValidLesson(lesson lesson_domain.Input) error {
	if lesson.Name == "" {
		return errors.New("name lesson cannot be empty")
	}

	if lesson.CourseID.Hex() == "" || lesson.CourseID.IsZero() {
		return errors.New("name lesson cannot be empty or data null")
	}

	if lesson.Level == 0 {
		return errors.New("level cannot be empty")
	}

	if lesson.Content == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}

func IsValidUnit(unit unit_domain.Input) error {
	if unit.LessonID.IsZero() {
		return errors.New("lesson name cannot be empty")
	}

	if unit.Name == "" {
		return errors.New("name cannot be empty")
	}

	if unit.Content == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}
