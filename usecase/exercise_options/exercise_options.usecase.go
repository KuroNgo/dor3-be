package exercise_options_usecase

import (
	exercise_options_domain "clean-architecture/domain/exercise_options"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type exerciseOptionsUseCase struct {
	exerciseOptionsRepository exercise_options_domain.IExerciseOptionRepository
	contextTimeout            time.Duration
}

func NewExerciseOptionsUseCase(exerciseOptionsRepository exercise_options_domain.IExerciseOptionRepository, timeout time.Duration) exercise_options_domain.IExerciseOptionUseCase {
	return &exerciseOptionsUseCase{
		exerciseOptionsRepository: exerciseOptionsRepository,
		contextTimeout:            timeout,
	}
}

func (e *exerciseOptionsUseCase) FetchManyByQuestionID(ctx context.Context, questionID string) (exercise_options_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseOptionsRepository.FetchManyByQuestionID(ctx, questionID)
	if err != nil {
		return exercise_options_domain.Response{}, err
	}

	return data, nil
}

func (e *exerciseOptionsUseCase) UpdateOne(ctx context.Context, exerciseOptions *exercise_options_domain.ExerciseOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseOptionsRepository.UpdateOne(ctx, exerciseOptions)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *exerciseOptionsUseCase) CreateOne(ctx context.Context, exerciseOptions *exercise_options_domain.ExerciseOptions) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseOptionsRepository.CreateOne(ctx, exerciseOptions)
	if err != nil {
		return err
	}

	return nil
}

func (e *exerciseOptionsUseCase) DeleteOne(ctx context.Context, optionsID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseOptionsRepository.DeleteOne(ctx, optionsID)
	if err != nil {
		return err
	}

	return nil
}
