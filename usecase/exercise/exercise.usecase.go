package exercise_usecase

import (
	exercise_domain "clean-architecture/domain/exercise"
	"context"
	"time"
)

type exerciseUseCase struct {
	exerciseRepository exercise_domain.IExerciseRepository
	contextTimeout     time.Duration
}

func NewExerciseUseCase(exerciseRepository exercise_domain.IExerciseRepository, timeout time.Duration) exercise_domain.IExerciseUseCase {
	return &exerciseUseCase{
		exerciseRepository: exerciseRepository,
		contextTimeout:     timeout,
	}
}

func (e *exerciseUseCase) FetchMany(ctx context.Context, page string) (exercise_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	vocabulary, err := e.exerciseRepository.FetchMany(ctx, page)
	if err != nil {
		return exercise_domain.Response{}, err
	}

	return vocabulary, err
}

func (e *exerciseUseCase) UpdateOne(ctx context.Context, exerciseID string, exercise exercise_domain.Exercise) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseRepository.UpdateOne(ctx, exerciseID, exercise)
	if err != nil {
		return err
	}

	return nil
}

func (e *exerciseUseCase) CreateOne(ctx context.Context, exercise *exercise_domain.Exercise) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseRepository.CreateOne(ctx, exercise)
	if err != nil {
		return err
	}

	return nil
}

func (e *exerciseUseCase) UpsertOne(ctx context.Context, id string, exercise *exercise_domain.Exercise) (exercise_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	exercises, err := e.exerciseRepository.UpsertOne(ctx, id, exercise)
	if err != nil {
		return exercise_domain.Response{}, err
	}

	return exercises, nil
}

func (e *exerciseUseCase) DeleteOne(ctx context.Context, exerciseID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseRepository.DeleteOne(ctx, exerciseID)
	if err != nil {
		return err
	}

	return err
}
