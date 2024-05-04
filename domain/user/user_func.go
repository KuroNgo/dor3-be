package user_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
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
	Create(ctx context.Context, user *User) error
	UpdateVerify(ctx context.Context, user *User) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, userID string) error
	Login(ctx context.Context, request SignIn) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	UpdateImage(ctx context.Context, userID string, imageURL string) error
	GetByID(ctx context.Context, id string) (*User, error)
}
