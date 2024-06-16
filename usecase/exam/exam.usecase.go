package exam_usecase

import (
	exam_domain "clean-architecture/domain/exam"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type examUseCase struct {
	examRepository exam_domain.IExamRepository
	contextTimeout time.Duration
}

func (e *examUseCase) FetchOneByUnitIDInUser(ctx context.Context, userID primitive.ObjectID, unitID string) (exam_domain.Exam, error) {
	//TODO implement me
	panic("implement me")
}

func NewExamUseCase(examRepository exam_domain.IExamRepository, timeout time.Duration) exam_domain.IExamUseCase {
	return &examUseCase{
		examRepository: examRepository,
		contextTimeout: timeout,
	}
}

func (e *examUseCase) FetchExamByIDInAdmin(ctx context.Context, id string) (exam_domain.Exam, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examRepository.FetchExamByIDInAdmin(ctx, id)
	if err != nil {
		return exam_domain.Exam{}, err
	}

	return data, nil
}

func (e *examUseCase) FetchManyInAdmin(ctx context.Context, page string) ([]exam_domain.Exam, exam_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, detail, err := e.examRepository.FetchManyInAdmin(ctx, page)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	return data, detail, nil
}

func (e *examUseCase) FetchOneByUnitIDInAdmin(ctx context.Context, unitID string) (exam_domain.Exam, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examRepository.FetchOneByUnitIDInAdmin(ctx, unitID)
	if err != nil {
		return exam_domain.Exam{}, err
	}

	return data, nil
}

func (e *examUseCase) FetchManyByUnitIDInAdmin(ctx context.Context, unitID string, page string) ([]exam_domain.Exam, exam_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, detail, err := e.examRepository.FetchManyByUnitIDInAdmin(ctx, unitID, page)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	return data, detail, nil
}

func (e *examUseCase) CreateOneInAdmin(ctx context.Context, exam *exam_domain.Exam) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examRepository.CreateOneInAdmin(ctx, exam)
	if err != nil {
		return err
	}

	return nil
}

func (e *examUseCase) UpdateOneInAdmin(ctx context.Context, exam *exam_domain.Exam) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examRepository.UpdateOneInAdmin(ctx, exam)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *examUseCase) DeleteOneInAdmin(ctx context.Context, examID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examRepository.DeleteOneInAdmin(ctx, examID)
	if err != nil {
		return err
	}

	return nil
}

func (e *examUseCase) UpdateCompletedInUser(ctx context.Context, exam *exam_domain.Exam) error {
	//TODO implement me
	panic("implement me")
}
