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

func (f *feedbackUseCase) FetchManyInAdmin(ctx context.Context, page string) (feedback_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	feedback, err := f.feedbackRepository.FetchManyInAdmin(ctx, page)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	return feedback, err
}

func (f *feedbackUseCase) FetchByUserIDInAdmin(ctx context.Context, userID string, page string) (feedback_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	feedback, err := f.feedbackRepository.FetchByUserIDInAdmin(ctx, userID, page)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	return feedback, err
}

func (f *feedbackUseCase) FetchBySubmittedDateInAdmin(ctx context.Context, date string, page string) (feedback_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	feedback, err := f.feedbackRepository.FetchBySubmittedDateInAdmin(ctx, date, page)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	return feedback, err
}

func (f *feedbackUseCase) CreateOneInUser(ctx context.Context, feedback *feedback_domain.Feedback) error {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	err := f.feedbackRepository.CreateOneInUser(ctx, feedback)
	if err != nil {
		return err
	}

	return nil
}

func (f *feedbackUseCase) DeleteOneInAdmin(ctx context.Context, feedbackID string) error {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	err := f.feedbackRepository.DeleteOneInAdmin(ctx, feedbackID)
	if err != nil {
		return err
	}

	return nil
}

func (f *feedbackUseCase) UpdateSeenInAdmin(ctx context.Context, id string, isSeen int) error {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	err := f.feedbackRepository.UpdateSeenInAdmin(ctx, id, isSeen)
	if err != nil {
		return err
	}

	return nil
}
