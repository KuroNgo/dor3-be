package exercise_question_usecase

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

func (e *exerciseQuestionUseCase) FetchOneByExerciseIDInAdmin(ctx context.Context, exerciseID string) (exercise_questions_domain.ExerciseQuestionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	vocabulary, err := e.exerciseQuestionRepository.FetchOneByExerciseIDInAdmin(ctx, exerciseID)
	if err != nil {
		return exercise_questions_domain.ExerciseQuestionResponse{}, err
	}

	return vocabulary, err

}

func NewExerciseQuestionUseCase(exerciseQuestionRepository exercise_questions_domain.IExerciseQuestionRepository, timeout time.Duration) exercise_questions_domain.IExerciseQuestionUseCase {
	return &exerciseQuestionUseCase{
		exerciseQuestionRepository: exerciseQuestionRepository,
		contextTimeout:             timeout,
	}
}

func (e *exerciseQuestionUseCase) FetchManyInAdmin(ctx context.Context, page string) (exercise_questions_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.FetchManyInAdmin(ctx, page)
	if err != nil {
		return exercise_questions_domain.Response{}, err
	}

	return data, nil
}

func (e *exerciseQuestionUseCase) FetchByIDInAdmin(ctx context.Context, id string) (exercise_questions_domain.ExerciseQuestionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	vocabulary, err := e.exerciseQuestionRepository.FetchByIDInAdmin(ctx, id)
	if err != nil {
		return exercise_questions_domain.ExerciseQuestionResponse{}, err
	}

	return vocabulary, err
}

func (e *exerciseQuestionUseCase) FetchManyByExerciseIDInAdmin(ctx context.Context, exerciseID string) (exercise_questions_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.FetchManyByExerciseIDInAdmin(ctx, exerciseID)
	if err != nil {
		return exercise_questions_domain.Response{}, err
	}

	return data, nil
}

func (e *exerciseQuestionUseCase) UpdateOneInAdmin(ctx context.Context, exerciseQuestion *exercise_questions_domain.ExerciseQuestion) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.exerciseQuestionRepository.UpdateOneInAdmin(ctx, exerciseQuestion)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *exerciseQuestionUseCase) CreateOneInAdmin(ctx context.Context, exerciseQuestion *exercise_questions_domain.ExerciseQuestion) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseQuestionRepository.CreateOneInAdmin(ctx, exerciseQuestion)
	if err != nil {
		return err
	}

	return nil
}

func (e *exerciseQuestionUseCase) DeleteOneInAdmin(ctx context.Context, exerciseID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.exerciseQuestionRepository.DeleteOneInAdmin(ctx, exerciseID)
	if err != nil {
		return err
	}

	return nil
}
