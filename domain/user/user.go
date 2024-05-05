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
	User []User `bson:"user" json:"user"`
}

type Statistics struct {
	Total int64 `bson:"total" json:"total"`
}

//go:generate mockery --name IUserRepository
type IUserRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)

	DeleteOne(ctx context.Context, userID string) error
	Login(ctx context.Context, request SignIn) (*User, error)
	Create(ctx context.Context, user *User) error

	Update(ctx context.Context, user *User) error
	UpdateVerify(ctx context.Context, user *User) (*mongo.UpdateResult, error)
	UpsertOne(ctx context.Context, email string, user *User) (*User, error)
	UpdateImage(ctx context.Context, userID string, imageURL string) error
}
