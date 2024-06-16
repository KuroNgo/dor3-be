package exam_result_usecase

import (
	exam_result_domain "clean-architecture/domain/exam_result"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (e *examResultUseCase) FetchManyInUser(ctx context.Context, examID string, userID primitive.ObjectID) (exam_result_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examResultRepository.FetchManyInUser(ctx, examID)
	if err != nil {
		return exam_result_domain.Response{}, err
	}

	return data, nil
}

func (e *examResultUseCase) GetResultByIDInUser(ctx context.Context, userID string) (exam_result_domain.ExamResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examResultRepository.GetResultByIDInUser(ctx, userID)
	if err != nil {
		return exam_result_domain.ExamResult{}, err
	}

	return data, nil
}

func (e *examResultUseCase) FetchManyByExamIDInUser(ctx context.Context, examID string) (exam_result_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examResultRepository.FetchManyInUser(ctx, examID)
	if err != nil {
		return exam_result_domain.Response{}, err
	}

	return data, nil

}

func (e *examResultUseCase) GetResultsByExamIDInUser(ctx context.Context, userID string, examID string) (exam_result_domain.ExamResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examResultRepository.GetResultsByExamIDInUser(ctx, userID, examID)
	if err != nil {
		return exam_result_domain.ExamResult{}, err
	}

	return data, nil
}

func (e *examResultUseCase) CreateOneInUser(ctx context.Context, examResult *exam_result_domain.ExamResult) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examResultRepository.CreateOneInUser(ctx, examResult)
	if err != nil {
		return err
	}

	return nil
}

func (e *examResultUseCase) UpdateStatusInUser(ctx context.Context, examResultID string, status int) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examResultRepository.UpdateStatusInUser(ctx, examResultID, status)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *examResultUseCase) CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	value := e.examResultRepository.CalculateScore(ctx, correctAnswers, totalQuestions)
	return value
}

func (e *examResultUseCase) CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64 {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	value := e.examResultRepository.CalculatePercentage(ctx, correctAnswers, totalQuestions)
	return value
}

func (e *examResultUseCase) DeleteOneInUser(ctx context.Context, examResultID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examResultRepository.DeleteOneInUser(ctx, examResultID)
	if err != nil {
		return err
	}

	return nil
}
