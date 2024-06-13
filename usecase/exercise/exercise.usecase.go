package exercise_usecase

import (
	exercise_domain "clean-architecture/domain/exercise"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
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

func (e *exerciseUseCase) FetchByIDInAdmin(ctx context.Context, id string) (exercise_domain.Exercise, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	vocabulary, err := e.exerciseRepository.FetchByIDInAdmin(ctx, id)
	if err != nil {
		return exercise_domain.Exercise{}, err
	}

	return vocabulary, err
}

func (e *exerciseUseCase) FetchOneByUnitIDInAdmin(ctx context.Context, unitID string) (exercise_domain.Exercise, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	vocabulary, err := e.exerciseRepository.FetchOneByUnitIDInAdmin(ctx, unitID)
	if err != nil {
		return exercise_domain.Exercise{}, err
	}

	return vocabulary, err
}

func (e *exerciseUseCase) FetchManyByUnitIDInAdmin(ctx context.Context, unitID string, page string) ([]exercise_domain.Exercise, exercise_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	vocabulary, detail, err := e.exerciseRepository.FetchManyByUnitIDInAdmin(ctx, unitID, page)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	return vocabulary, detail, err
}

func (e *exerciseUseCase) FetchManyInAdmin(ctx context.Context, page string) ([]exercise_domain.Exercise, exercise_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	exercises, detail, err := e.exerciseRepository.FetchManyInAdmin(ctx, page)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	return exercises, detail, err
}

func (e *exerciseUseCase) UpdateOneInAdmin(ctx context.Context, exercise *exercise_domain.Exercise) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseRepository.UpdateOneInAdmin(ctx, exercise)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *exerciseUseCase) CreateOneInAdmin(ctx context.Context, exercise *exercise_domain.Exercise) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseRepository.CreateOneInAdmin(ctx, exercise)
	if err != nil {
		return err
	}

	return nil
}

func (e *exerciseUseCase) DeleteOneInAdmin(ctx context.Context, exerciseID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseRepository.DeleteOneInAdmin(ctx, exerciseID)
	if err != nil {
		return err
	}

	return err
}
