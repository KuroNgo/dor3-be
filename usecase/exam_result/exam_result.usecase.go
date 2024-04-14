package exam_result_usecase

import (
	exam_result_domain "clean-architecture/domain/exam_result"
	"context"
	"time"
)

type examResultUseCase struct {
	examResultRepository exam_result_domain.IExamResultRepository
	contextTimeout       time.Duration
}

func NewExamResultUseCase(examResultRepository exam_result_domain.IExamResultRepository, timeout time.Duration) exam_result_domain.IExamResultUseCase {
	return &examResultUseCase{
		examResultRepository: examResultRepository,
		contextTimeout:       timeout,
	}
}

func (e *examResultUseCase) FetchMany(ctx context.Context, page string) (exam_result_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examResultRepository.FetchMany(ctx, page)
	if err != nil {
		return exam_result_domain.Response{}, err
	}

	return data, nil
}

func (e *examResultUseCase) FetchManyByExamID(ctx context.Context, examID string) (exam_result_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examResultRepository.FetchManyByExamID(ctx, examID)
	if err != nil {
		return exam_result_domain.Response{}, err
	}

	return data, nil

}

func (e *examResultUseCase) CreateOne(ctx context.Context, examResult *exam_result_domain.ExamResult) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examResultRepository.CreateOne(ctx, examResult)
	if err != nil {
		return err
	}

	return nil
}

func (e *examResultUseCase) DeleteOne(ctx context.Context, examResultID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examResultRepository.DeleteOne(ctx, examResultID)
	if err != nil {
		return err
	}

	return nil
}
