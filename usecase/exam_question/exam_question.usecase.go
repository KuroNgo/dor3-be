package exam_question_usecase

import (
	exam_question_domain "clean-architecture/domain/exam_question"
	"context"
	"time"
)

type examQuestionUseCase struct {
	examQuestionRepository exam_question_domain.IExamQuestionRepository
	contextTimeout         time.Duration
}

func (e *examQuestionUseCase) FetchMany(ctx context.Context, page string) (exam_question_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examQuestionRepository.FetchMany(ctx, page)
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	return data, nil
}

func (e *examQuestionUseCase) FetchManyByExamID(ctx context.Context, examID string) (exam_question_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examQuestionRepository.FetchManyByExamID(ctx, examID)
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	return data, nil
}

func (e *examQuestionUseCase) UpdateOne(ctx context.Context, examQuestionID string, examQuestion exam_question_domain.ExamQuestion) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examQuestionRepository.UpdateOne(ctx, examQuestionID, examQuestion)
	if err != nil {
		return err
	}

	return nil
}

func (e *examQuestionUseCase) CreateOne(ctx context.Context, examQuestion *exam_question_domain.ExamQuestion) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examQuestionRepository.CreateOne(ctx, examQuestion)
	if err != nil {
		return err
	}

	return nil
}

func (e *examQuestionUseCase) DeleteOne(ctx context.Context, examID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examQuestionRepository.DeleteOne(ctx, examID)
	if err != nil {
		return err
	}

	return nil
}

func NewExamQuestionUseCase(examQuestionRepository exam_question_domain.IExamQuestionRepository, timeout time.Duration) exam_question_domain.IExamQuestionUseCase {
	return &examQuestionUseCase{
		examQuestionRepository: examQuestionRepository,
		contextTimeout:         timeout,
	}
}
