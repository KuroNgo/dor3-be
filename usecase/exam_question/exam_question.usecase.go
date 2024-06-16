package exam_question_usecase

import (
	exam_question_domain "clean-architecture/domain/exam_question"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type examQuestionUseCase struct {
	examQuestionRepository exam_question_domain.IExamQuestionRepository
	contextTimeout         time.Duration
}

func NewExamQuestionUseCase(examQuestionRepository exam_question_domain.IExamQuestionRepository, timeout time.Duration) exam_question_domain.IExamQuestionUseCase {
	return &examQuestionUseCase{
		examQuestionRepository: examQuestionRepository,
		contextTimeout:         timeout,
	}
}

func (e *examQuestionUseCase) FetchManyInAdmin(ctx context.Context, page string) (exam_question_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examQuestionRepository.FetchManyInAdmin(ctx, page)
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	return data, nil
}

func (e *examQuestionUseCase) FetchOneByExamIDInAdmin(ctx context.Context, examID string) (exam_question_domain.ExamQuestionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examQuestionRepository.FetchOneByExamIDInAdmin(ctx, examID)
	if err != nil {
		return exam_question_domain.ExamQuestionResponse{}, err
	}

	return data, nil
}

func (e *examQuestionUseCase) FetchQuestionByIDInAdmin(ctx context.Context, id string) (exam_question_domain.ExamQuestion, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examQuestionRepository.FetchQuestionByIDInAdmin(ctx, id)
	if err != nil {
		return exam_question_domain.ExamQuestion{}, err
	}

	return data, nil
}

func (e *examQuestionUseCase) FetchManyByExamIDInAdmin(ctx context.Context, examID string, page string) (exam_question_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examQuestionRepository.FetchManyByExamIDInAdmin(ctx, examID, page)
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	return data, nil
}

func (e *examQuestionUseCase) UpdateOneInAdmin(ctx context.Context, examQuestion *exam_question_domain.ExamQuestion) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	data, err := e.examQuestionRepository.UpdateOneInAdmin(ctx, examQuestion)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e *examQuestionUseCase) CreateOneInAdmin(ctx context.Context, examQuestion *exam_question_domain.ExamQuestion) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examQuestionRepository.CreateOneInAdmin(ctx, examQuestion)
	if err != nil {
		return err
	}

	return nil
}

func (e *examQuestionUseCase) DeleteOneInAdmin(ctx context.Context, examID string) error {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	err := e.examQuestionRepository.DeleteOneInAdmin(ctx, examID)
	if err != nil {
		return err
	}

	return nil
}
