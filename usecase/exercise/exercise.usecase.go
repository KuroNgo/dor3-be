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

func (e *exerciseUseCase) FetchOneByUnitID(ctx context.Context, unitID string) (exercise_domain.ExerciseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewExerciseUseCase(exerciseRepository exercise_domain.IExerciseRepository, timeout time.Duration) exercise_domain.IExerciseUseCase {
	return &exerciseUseCase{
		exerciseRepository: exerciseRepository,
		contextTimeout:     timeout,
	}
}

func (e *exerciseUseCase) FetchManyByLessonID(ctx context.Context, unitID string) ([]exercise_domain.ExerciseResponse, exercise_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	vocabulary, detail, err := e.exerciseRepository.FetchManyByLessonID(ctx, unitID)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	return vocabulary, detail, err
}

func (e *exerciseUseCase) FetchManyByUnitID(ctx context.Context, unitID string) ([]exercise_domain.ExerciseResponse, exercise_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	vocabulary, detail, err := e.exerciseRepository.FetchManyByUnitID(ctx, unitID)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	return vocabulary, detail, err
}

func (e *exerciseUseCase) UpdateCompleted(ctx context.Context, exerciseID string, isComplete int) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseRepository.UpdateCompleted(ctx, exerciseID, isComplete)
	if err != nil {
		return err
	}

	return nil
}

func (e *exerciseUseCase) FetchMany(ctx context.Context, page string) ([]exercise_domain.ExerciseResponse, exercise_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	exercises, detail, err := e.exerciseRepository.FetchMany(ctx, page)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	return exercises, detail, err
}

func (e *exerciseUseCase) UpdateOne(ctx context.Context, exercise *exercise_domain.Exercise) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseRepository.UpdateOne(ctx, exercise)
	if err != nil {
		return nil, err
	}
	return data, nil
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

func (e *exerciseUseCase) DeleteOne(ctx context.Context, exerciseID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseRepository.DeleteOne(ctx, exerciseID)
	if err != nil {
		return err
	}

	return err
}
