package user_attempt_domain

import "context"

type IUserProcessUseCase interface {
	FetchManyByUserID(c context.Context) (Response, error)
	CreateOneByUserID(c context.Context, userID string) error
	DeleteOneByUserID(c context.Context, userID string) error
}