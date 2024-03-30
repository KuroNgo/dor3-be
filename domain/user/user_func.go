package user_domain

import (
	"context"
)

type SignUp struct {
	FullName   string `json:"full_name"  bson:"full_name"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	AvatarURL  string `json:"avatar_url"  bson:"avatar_url"`
	Specialize string `json:"specialize"  bson:"specialize"`
	Phone      string `json:"phone" bson:"phone"`
}

type SignIn struct {
	Email    string `json:"email" bson:"email"`
	Password string `bson:"password"  json:"password"`
}

type Input struct {
	FullName   string `bson:"full_name"  json:"full_name"`
	Email      string `bson:"email"  json:"email"`
	Password   string `bson:"password"  json:"password"`
	AvatarURL  string `bson:"avatar_url"  json:"avatar_url"`
	Specialize string `bson:"specialize"  json:"specialize"`
	Phone      string `bson:"phone"   json:"phone"`
}

//go:generate mockery --name IUserUseCase
type IUserUseCase interface {
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, userID string, user User) error
	Delete(ctx context.Context, userID string, user User) error
	Login(c context.Context, request SignIn) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdateImage(c context.Context, userID string, imageURL string) error
	GetByID(ctx context.Context, id string) (*User, error)
}
