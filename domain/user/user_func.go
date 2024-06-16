package user_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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

type UserInput struct {
	ID        primitive.ObjectID `bson:"_id"  json:"_id"`
	FullName  string             `bson:"full_name"  json:"full_name"`
	Email     string             `bson:"email"  json:"email"`
	Role      string             `bson:"role" json:"role"`
	AvatarURL string             `bson:"avatar_url"  json:"avatar_url"`
	Phone     string             `bson:"phone"   json:"phone"`
	Verified  bool               `bson:"verify"   json:"verify"`
	Provider  string             `bson:"provider" json:"provider"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

//go:generate mockery --name IUserUseCase
type IUserUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	CheckVerify(ctx context.Context, verificationCode string) bool
	GetByVerificationCode(ctx context.Context, verificationCode string) (*User, error)

	Create(ctx context.Context, user *User) error
	Login(ctx context.Context, request SignIn) (*User, error)
	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, user *User) error
	UpdateVerify(ctx context.Context, user *User) (*mongo.UpdateResult, error)
	UpdateVerifyForChangePassword(ctx context.Context, user *User) (*mongo.UpdateResult, error)
	UpdateImage(ctx context.Context, userID string, imageURL string) error
	Delete(ctx context.Context, userID string) error
	UniqueVerificationCode(ctx context.Context, verificationCode string) bool
}
