package activity_usecase

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	"context"
	"time"
)

type activityUseCaseV2 struct {
	activityRepository activity_log_domain.IActivityRepositoryV2
	contextTimeout     time.Duration
}

func NewActivityUseCaseV2(activityRepository activity_log_domain.IActivityRepositoryV2, timeout time.Duration) activity_log_domain.IActivityUseCaseV2 {
	return &activityUseCaseV2{
		activityRepository: activityRepository,
		contextTimeout:     timeout,
	}
}

func (a *activityUseCaseV2) CreateOne(ctx context.Context, log activity_log_domain.ActivityLog) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.activityRepository.CreateOne(ctx, log)
	if err != nil {
		return err
	}

	return nil
}
