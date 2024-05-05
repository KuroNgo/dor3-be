package exam_usecase

import (
	exam_domain "clean-architecture/domain/exam"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type examUseCase struct {
	examRepository exam_domain.IExamRepository
	contextTimeout time.Duration
}

func (e *examUseCase) FetchOneByUnitID(ctx context.Context, unitID string) (exam_domain.ExamResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewExamUseCase(examRepository exam_domain.IExamRepository, timeout time.Duration) exam_domain.IExamUseCase {
	return &examUseCase{
		examRepository: examRepository,
		contextTimeout: timeout,
	}
}

func (e *examUseCase) FetchMany(ctx context.Context, page string) ([]exam_domain.ExamResponse, exam_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, detail, err := e.examRepository.FetchMany(ctx, page)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	return data, detail, nil
}

func (e *examUseCase) FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]exam_domain.ExamResponse, exam_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, detail, err := e.examRepository.FetchManyByUnitID(ctx, unitID, page)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	return data, detail, nil
}

func (e *examUseCase) UpdateOne(ctx context.Context, exam *exam_domain.Exam) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examRepository.UpdateOne(ctx, exam)
	if err != nil {
		return nil, err
	}

	return data, nil
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

func (e *examUseCase) DeleteOne(ctx context.Context, examID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examRepository.DeleteOne(ctx, examID)
	if err != nil {
		return err
	}

	return nil
}
