package activity_log_domain

import (
	"context"
)

type IActivityUseCase interface {
	CreateOne(ctx context.Context, log ActivityLog) error
	DeleteOne(ctx context.Context) error
	FetchMany(ctx context.Context) ([]ActivityLog, error)
	FetchByUserName(ctx context.Context, username string) (ActivityLog, error)
}
