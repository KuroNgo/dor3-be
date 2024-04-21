package exam_answer_usecase

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	"context"
	"time"
)

type examAnswerUseCase struct {
	examAnswerRepository exam_answer_domain.IExamAnswerRepository
	contextTimeout       time.Duration
}

func NewExamAnswerUseCase(examAnswerRepository exam_answer_domain.IExamAnswerRepository, timeout time.Duration) exam_answer_domain.IExamAnswerUseCase {
	return &examAnswerUseCase{
		examAnswerRepository: examAnswerRepository,
		contextTimeout:       timeout,
	}
}

func (e *examAnswerUseCase) DeleteAllAnswerByExamID(ctx context.Context, examID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examAnswerRepository.DeleteAllAnswerByExamID(ctx, examID)
	if err != nil {
		return err
	}

	return nil
}

func (e *examAnswerUseCase) FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (exam_answer_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examAnswerRepository.FetchManyAnswerByUserIDAndQuestionID(ctx, questionID, userID)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	return data, nil
}

func (e *examAnswerUseCase) CreateOne(ctx context.Context, examAnswer *exam_answer_domain.ExamAnswer) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examAnswerRepository.CreateOne(ctx, examAnswer)
	if err != nil {
		return err
	}

	return nil
}

func (e *examAnswerUseCase) DeleteOne(ctx context.Context, examID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examAnswerRepository.DeleteOne(ctx, examID)
	if err != nil {
		return err
	}

	return nil
}
