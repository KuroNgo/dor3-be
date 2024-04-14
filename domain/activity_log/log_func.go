package activity_log_domain

import (
	"context"
)

type IActivityUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	CreateOne(ctx context.Context, log ActivityLog) error
}
