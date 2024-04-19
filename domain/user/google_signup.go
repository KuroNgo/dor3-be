package user_domain

import (
	"context"
)

//go:generate mockery --name IGoogleAuthUseCase
type IGoogleAuthUseCase interface {
	UpsertUser(c context.Context, email string, user *User) (*User, error)
}
