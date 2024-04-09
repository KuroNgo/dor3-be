package activity_log_domain

import (
	"context"
	"time"
)

type IActivityUseCase interface {
	CreateOne(ctx context.Context, log ActivityLog) error
	DeleteOne(ctx context.Context, logID string) error
	DeleteOneByTime(ctx context.Context, time time.Duration) error
	FetchMany(ctx context.Context) ([]ActivityLog, error)
	FetchByUserName(ctx context.Context, username string) (ActivityLog, error)
}
