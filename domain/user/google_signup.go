package user_domain

import (
	"context"
)

//go:generate mockery --name IGoogleAuthUseCase
type IGoogleAuthUseCase interface {
	UpsertUser(ctx context.Context, email string, user *User) (*User, error)
}
