package exercise_result

import (
	exercise_result_domain "clean-architecture/domain/exercise_result"
	"context"
	"time"
)

type exerciseResultUseCase struct {
	exerciseQuestionRepository exercise_result_domain.IExerciseResultRepository
	contextTimeout             time.Duration
}

func (e *exerciseResultUseCase) FetchMany(ctx context.Context, page string) (exercise_result_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.FetchMany(ctx, page)
	if err != nil {
		return exercise_result_domain.Response{}, err
	}

	return data, nil
}

func (e *exerciseResultUseCase) FetchManyByExamID(ctx context.Context, exerciseID string) (exercise_result_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.FetchManyByExamID(ctx, exerciseID)
	if err != nil {
		return exercise_result_domain.Response{}, err
	}

	return data, nil
}

func (e *exerciseResultUseCase) CreateOne(ctx context.Context, exerciseResult *exercise_result_domain.ExerciseResult) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseQuestionRepository.CreateOne(ctx, exerciseResult)
	if err != nil {
		return err
	}

	return nil
}

func (e *exerciseResultUseCase) DeleteOne(ctx context.Context, exerciseResultID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseQuestionRepository.DeleteOne(ctx, exerciseResultID)
	if err != nil {
		return err
	}

	return nil
}

func NewExerciseQuestionUseCase(exerciseQuestionRepository exercise_result_domain.IExerciseResultRepository, timeout time.Duration) exercise_result_domain.IExerciseResultUseCase {
	return &exerciseResultUseCase{
		exerciseQuestionRepository: exerciseQuestionRepository,
		contextTimeout:             timeout,
	}
}
