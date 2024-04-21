package exercise_question

import (
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type exerciseQuestionUseCase struct {
	exerciseQuestionRepository exercise_questions_domain.IExerciseQuestionRepository
	contextTimeout             time.Duration
}

func NewExerciseQuestionUseCase(exerciseQuestionRepository exercise_questions_domain.IExerciseQuestionRepository, timeout time.Duration) exercise_questions_domain.IExerciseQuestionUseCase {
	return &exerciseQuestionUseCase{
		exerciseQuestionRepository: exerciseQuestionRepository,
		contextTimeout:             timeout,
	}
}

func (e *exerciseQuestionUseCase) FetchMany(ctx context.Context, page string) (exercise_questions_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.FetchMany(ctx, page)
	if err != nil {
		return exercise_questions_domain.Response{}, err
	}

	return data, nil
}

func (e *exerciseQuestionUseCase) FetchManyByExerciseID(ctx context.Context, exerciseID string) (exercise_questions_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.FetchManyByExamID(ctx, exerciseID)
	if err != nil {
		return exercise_questions_domain.Response{}, err
	}

	return data, nil
}

func (e *exerciseQuestionUseCase) UpdateOne(ctx context.Context, exerciseQuestion *exercise_questions_domain.ExerciseQuestion) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.UpdateOne(ctx, exerciseQuestion)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *exerciseQuestionUseCase) CreateOne(ctx context.Context, exerciseQuestion *exercise_questions_domain.ExerciseQuestion) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseQuestionRepository.CreateOne(ctx, exerciseQuestion)
	if err != nil {
		return err
	}

	return nil
}

func (e *exerciseQuestionUseCase) DeleteOne(ctx context.Context, exerciseID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseQuestionRepository.DeleteOne(ctx, exerciseID)
	if err != nil {
		return err
	}

	return nil
}
