package exam_options_usecase

import (
	exam_options_domain "clean-architecture/domain/exam_options"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type examOptionsUseCase struct {
	examOptionsRepository exam_options_domain.IExamOptionRepository
	contextTimeout        time.Duration
}

func NewExamOptionsUseCase(examOptionsRepository exam_options_domain.IExamOptionRepository, timeout time.Duration) exam_options_domain.IExamOptionsUseCase {
	return &examOptionsUseCase{
		examOptionsRepository: examOptionsRepository,
		contextTimeout:        timeout,
	}
}

func (e *examOptionsUseCase) FetchManyByQuestionID(ctx context.Context, questionID string) (exam_options_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examOptionsRepository.FetchManyByQuestionID(ctx, questionID)
	if err != nil {
		return exam_options_domain.Response{}, err
	}

	return data, nil
}

func (e *examOptionsUseCase) UpdateOne(ctx context.Context, examOptions *exam_options_domain.ExamOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examOptionsRepository.UpdateOne(ctx, examOptions)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *examOptionsUseCase) CreateOne(ctx context.Context, examOptions *exam_options_domain.ExamOptions) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examOptionsRepository.CreateOne(ctx, examOptions)
	if err != nil {
		return err
	}

	return nil
}

func (e *examOptionsUseCase) DeleteOne(ctx context.Context, examID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examOptionsRepository.DeleteOne(ctx, examID)
	if err != nil {
		return err
	}

	return nil
}
