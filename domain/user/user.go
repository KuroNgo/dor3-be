package user_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id"`
	IP               string             `bson:"ip" json:"ip"`
	FullName         string             `bson:"full_name"  json:"full_name"`
	Email            string             `bson:"email"  json:"email"`
	Password         string             `bson:"password"  json:"password"`
	Role             string             `bson:"role" json:"role"`
	CoverURL         string             `bson:"cover_url" json:"cover_url"`
	AvatarURL        string             `bson:"avatar_url"  json:"avatar_url"`
	AssetID          string             `bson:"asset_id"  json:"asset_id"`
	Phone            string             `bson:"phone"   json:"phone"`
	Provider         string             `bson:"provider" json:"provider"`
	Verified         bool               `bson:"verified" json:"verified"`
	VerificationCode string             `bson:"verification_code" json:"verification_code"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type Response struct {
	Total int64  `bson:"total" json:"total"`
	User  []User `bson:"user" json:"user"`
}

//go:generate mockery --name IUserRepository
type IUserRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*User, error)
	CheckVerify(ctx context.Context, verificationCode string) bool
	DeleteOne(ctx context.Context, userID string) error
	Login(ctx context.Context, request SignIn) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, user *User) error
	UpdateVerify(ctx context.Context, user *User) (*mongo.UpdateResult, error)
	UpdateVerifyForChangePassword(ctx context.Context, user *User) (*mongo.UpdateResult, error)
	UpsertOne(ctx context.Context, email string, user *UserInput) (*User, error)
	UpdateImage(ctx context.Context, userID string, imageURL string) error
	UniqueVerificationCode(ctx context.Context, verificationCode string) bool
}
