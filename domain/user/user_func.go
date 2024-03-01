package user_domain

import (
	"context"
)

type Input struct {
	FullName   string `json:"full_name"  bson:"full_name"`
	Email      string `json:"email" bson:"email"`
	AvatarURL  string `json:"avatar_url"  bson:"avatar_url"`
	Specialize string `json:"specialize"  bson:"specialize"`
	Photo      string `json:"photo" bson:"photo"`
}

//go:generate mockery --name IUserUseCase
type IUserUseCase interface {
	Create(c context.Context, user Input) error
	GetByEmail(c context.Context, email string) (*User, error)
	GetByUsername(c context.Context, username string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
}
