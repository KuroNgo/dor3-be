package exercise_answer

import (
	exercise_answer_domain "clean-architecture/domain/exercise_answer"
	"context"
	"time"
)

type exerciseAnswerUseCase struct {
	exerciseAnswerRepository exercise_answer_domain.IExerciseAnswerRepository
	contextTimeout           time.Duration
}

func NewExerciseAnswerUseCase(exerciseAnswerRepository exercise_answer_domain.IExerciseAnswerRepository, timeout time.Duration) exercise_answer_domain.IExerciseAnswerUseCase {
	return &exerciseAnswerUseCase{
		exerciseAnswerRepository: exerciseAnswerRepository,
		contextTimeout:           timeout,
	}
}

func (e *exerciseAnswerUseCase) FetchManyByQuestionID(ctx context.Context, questionID string) (exercise_answer_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseAnswerRepository.FetchManyByQuestionID(ctx, questionID)
	if err != nil {
		return exercise_answer_domain.Response{}, err
	}

	return data, nil
}

func (e *exerciseAnswerUseCase) CreateOne(ctx context.Context, exerciseAnswer *exercise_answer_domain.ExerciseAnswer) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseAnswerRepository.CreateOne(ctx, exerciseAnswer)
	if err != nil {
		return err
	}

	return nil
}

func (e *exerciseAnswerUseCase) DeleteOne(ctx context.Context, exerciseID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseAnswerRepository.DeleteOne(ctx, exerciseID)
	if err != nil {
		return err
	}

	return nil
}
