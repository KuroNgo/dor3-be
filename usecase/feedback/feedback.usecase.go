package feedback_usecase

import (
	feedback_domain "clean-architecture/domain/feedback"
	"context"
	"time"
)

type feedbackUseCase struct {
	feedbackRepository feedback_domain.IFeedbackRepository
	contextTimeout     time.Duration
}

func NewFeedbackUseCase(feedbackRepository feedback_domain.IFeedbackRepository, timeout time.Duration) feedback_domain.IFeedbackUseCase {
	return &feedbackUseCase{
		feedbackRepository: feedbackRepository,
		contextTimeout:     timeout,
	}
}

func (f *feedbackUseCase) FetchMany(ctx context.Context, page string) (feedback_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	feedback, err := f.feedbackRepository.FetchMany(ctx, page)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	return feedback, err
}

func (f *feedbackUseCase) FetchByUserID(ctx context.Context, userID string, page string) (feedback_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	feedback, err := f.feedbackRepository.FetchByUserID(ctx, userID, page)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	return feedback, err
}

func (f *feedbackUseCase) FetchBySubmittedDate(ctx context.Context, date string, page string) (feedback_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	feedback, err := f.feedbackRepository.FetchBySubmittedDate(ctx, date, page)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	return feedback, err
}

func (f *feedbackUseCase) CreateOneByUser(ctx context.Context, feedback *feedback_domain.Feedback) error {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	err := f.feedbackRepository.CreateOneByUser(ctx, feedback)
	if err != nil {
		return err
	}

	return nil
}

func (f *feedbackUseCase) DeleteOneByAdmin(ctx context.Context, feedbackID string) error {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	err := f.feedbackRepository.DeleteOneByAdmin(ctx, feedbackID)
	if err != nil {
		return err
	}

	return nil
}
