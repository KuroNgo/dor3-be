package user_domain

import (
	"context"
)

type SignUp struct {
	FullName   string `json:"full_name"  bson:"full_name"`
	Nickname   string `json:"nickname"  bson:"nickname"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	AvatarURL  string `json:"avatar_url"  bson:"avatar_url"`
	Specialize string `json:"specialize"  bson:"specialize"`
	Photo      string `json:"photo" bson:"photo"`
}

type SignIn struct {
	Email    string `json:"email" bson:"email"`
	Password string `bson:"password"  json:"password"`
}

//go:generate mockery --name IUserUseCase
type IUserUseCase interface {
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, userID string, user User) error
	Delete(ctx context.Context, userID string, user User) error
	Login(c context.Context, email string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}
