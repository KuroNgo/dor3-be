package internal

import (
	course_domain "clean-architecture/domain/course"
	exercise_domain "clean-architecture/domain/exercise"
	lesson_domain "clean-architecture/domain/lesson"
	quiz_domain "clean-architecture/domain/quiz"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func IsValidVocabulary(vocabulary vocabulary_domain.Input) error {
	if vocabulary.Word == "" {
		return errors.New("word cannot be empty")
	}

	if vocabulary.PartOfSpeech == "" {
		return errors.New("part of speech cannot be empty")
	}

	if vocabulary.Pronunciation == "" {
		return errors.New("pronunciation cannot be empty")
	}

	if vocabulary.ExampleEng == "" {
		return errors.New("example cannot be empty")
	}
	if vocabulary.ExplainVie == "" {
		return errors.New("example cannot be empty")
	}
	if vocabulary.ExampleVie == "" {
		return errors.New("example cannot be empty")
	}
	if vocabulary.ExplainEng == "" {
		return errors.New("example cannot be empty")
	}

	if vocabulary.FieldOfIT == "" {
		return errors.New("field of IT cannot be empty")
	}

	if vocabulary.LinkURL == "" {
		return errors.New("link URL cannot be empty")
	}

	return nil
}

func IsValidExercise(exercise exercise_domain.Input) error {
	//if exercise.Options == nil {
	//	return errors.New("option cannot be empty")
	//}
	if exercise.VocabularyID == primitive.NilObjectID {
		return errors.New("option cannot be empty")
	}
	if exercise.Title == "" {
		return errors.New("option cannot be empty")
	}
	if exercise.Type == "" {
		return errors.New("option cannot be empty")
	}
	if exercise.CorrectAns == "" {
		return errors.New("option cannot be empty")
	}
	return nil
}
