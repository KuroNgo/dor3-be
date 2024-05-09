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

type VerificationCode struct {
	VerificationCode string `json:"verification_code" bson:"verification_code"`
}

type ForgetPassword struct {
	Email string `json:"email" bson:"email"`
}

type ChangePassword struct {
	Password        string `json:"password" bson:"password"`
	PasswordCompare string `json:"password_compare" bson:"password_compare"`
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
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	CheckVerify(ctx context.Context, verificationCode string) bool
	GetByVerificationCode(ctx context.Context, verificationCode string) (*User, error)

	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID string) error
	Login(ctx context.Context, request SignIn) (*User, error)

	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, user *User) error
	UpdateVerify(ctx context.Context, user *User) (*mongo.UpdateResult, error)
	UpdateVerifyForChangePassword(ctx context.Context, user *User) (*mongo.UpdateResult, error)
	UpdateImage(ctx context.Context, userID string, imageURL string) error

	UniqueVerificationCode(ctx context.Context, verificationCode string) bool
}
