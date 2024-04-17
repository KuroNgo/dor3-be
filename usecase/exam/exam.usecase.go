package exam_usecase

import (
	exam_domain "clean-architecture/domain/exam"
	"context"
	"time"
)

type examUseCase struct {
	examRepository exam_domain.IExamRepository
	contextTimeout time.Duration
}

func NewExamUseCase(examRepository exam_domain.IExamRepository, timeout time.Duration) exam_domain.IExamUseCase {
	return &examUseCase{
		examRepository: examRepository,
		contextTimeout: timeout,
	}
}

func (e *examUseCase) FetchMany(ctx context.Context, page string) (exam_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	data, err := e.examRepository.FetchMany(ctx, page)

	if err != nil {
		return exam_domain.Response{}, err
	}

	return data, nil
}

func (e *examUseCase) FetchManyByUnitID(ctx context.Context, unitID string) (exam_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	data, err := e.examRepository.FetchManyByUnitID(ctx, unitID)

	if err != nil {
		return exam_domain.Response{}, err
	}

	return data, nil
}

func (e *examUseCase) UpdateOne(ctx context.Context, examID string, exam exam_domain.Exam) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	err := e.examRepository.UpdateOne(ctx, examID, exam)

	if err != nil {
		return err
	}

	return nil
}

func (e *examUseCase) CreateOne(ctx context.Context, exam *exam_domain.Exam) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	err := e.examRepository.CreateOne(ctx, exam)

	if err != nil {
		return err
	}

	return nil
}

func (e *examUseCase) UpdateCompleted(ctx context.Context, examID string, isComplete int) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	err := e.examRepository.UpdateCompleted(ctx, examID, isComplete)

	if err != nil {
		return err
	}

	return nil
}

func (e *examUseCase) DeleteOne(ctx context.Context, examID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	err := e.examRepository.DeleteOne(ctx, examID)

	if err != nil {
		return err
	}

	return nil
}