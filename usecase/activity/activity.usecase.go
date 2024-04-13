package activity_usecase

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	"context"
	"time"
)

type activityUseCase struct {
	activityRepository activity_log_domain.IActivityRepository
	contextTimeout     time.Duration
}

func (a *activityUseCase) DeleteOneByTime(ctx context.Context, time time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.activityRepository.DeleteOneByTime(ctx, time)
	if err != nil {
		return err
	}

	return nil
}

func NewActivityUseCase(activityRepository activity_log_domain.IActivityRepository, timeout time.Duration) activity_log_domain.IActivityUseCase {
	return &activityUseCase{
		activityRepository: activityRepository,
		contextTimeout:     timeout,
	}
}

func (a *activityUseCase) CreateOne(ctx context.Context, log activity_log_domain.ActivityLog) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.activityRepository.CreateOne(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

func (a *activityUseCase) DeleteOne(ctx context.Context, logID string) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.activityRepository.DeleteOne(ctx, logID)
	if err != nil {
		return err
	}

	return nil
}

func (a *activityUseCase) FetchMany(ctx context.Context, page string) (activity_log_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	log, err := a.activityRepository.FetchMany(ctx, page)
	if err != nil {
		return activity_log_domain.Response{}, err
	}

	return log, nil
}
