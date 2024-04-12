package activity_log_domain

import (
	"context"
	"time"
)

type IActivityUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	CreateOne(ctx context.Context, log ActivityLog) error
	DeleteOne(ctx context.Context, logID string) error
	DeleteOneByTime(ctx context.Context, time time.Duration) error
	FetchByUserName(ctx context.Context, username string) (Response, error)
}
