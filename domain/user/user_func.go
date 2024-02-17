package user_domain

import (
	"context"
)

//go:generate mockery --name IUserUseCase
type IUserUseCase interface {
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
}
