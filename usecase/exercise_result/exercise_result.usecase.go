package exercise_result_usecase

import (
	exercise_result_domain "clean-architecture/domain/exercise_result"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type exerciseResultUseCase struct {
	exerciseQuestionRepository exercise_result_domain.IExerciseResultRepository
	contextTimeout             time.Duration
}

func NewExerciseQuestionUseCase(exerciseQuestionRepository exercise_result_domain.IExerciseResultRepository, timeout time.Duration) exercise_result_domain.IExerciseResultUseCase {
	return &exerciseResultUseCase{
		exerciseQuestionRepository: exerciseQuestionRepository,
		contextTimeout:             timeout,
	}
}

func (e *exerciseResultUseCase) FetchManyByExamID(ctx context.Context, exerciseID string) (exercise_result_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.FetchManyByExerciseID(ctx, exerciseID)
	if err != nil {
		return exercise_result_domain.Response{}, err
	}

	return data, nil
}

func (e *exerciseResultUseCase) GetResultsByUserIDAndExamID(ctx context.Context, userID string, exerciseID string) (exercise_result_domain.ExerciseResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.GetResultsByUserIDAndExerciseID(ctx, userID, exerciseID)
	if err != nil {
		return exercise_result_domain.ExerciseResult{}, err
	}

	return data, nil
}

func (e *exerciseResultUseCase) UpdateStatus(ctx context.Context, exerciseResultID string, status int) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.UpdateStatus(ctx, exerciseResultID, status)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *exerciseResultUseCase) CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data := e.exerciseQuestionRepository.CalculateScore(ctx, correctAnswers, totalQuestions)

	return data
}

func (e *exerciseResultUseCase) CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64 {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data := e.exerciseQuestionRepository.CalculatePercentage(ctx, correctAnswers, totalQuestions)

	return data
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
