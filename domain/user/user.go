package user_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	FullName   string             `bson:"full_name"  json:"full_name"`
	Email      string             `bson:"email"  json:"email"`
	Password   string             `bson:"password"  json:"password"`
	AvatarURL  string             `bson:"avatar_url"  json:"avatar_url"`
	AssetID    string             `bson:"asset_id"  json:"asset_id"`
	Specialize string             `bson:"specialize"  json:"specialize"`
	Phone      string             `bson:"phone"   json:"phone"`
	Provider   string             `json:"provider" bson:"provider"`
	Verified   bool               `json:"verified" bson:"verified"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	Role       string             `bson:"role" json:"role"`
}

type Response struct {
	User  []User `bson:"user" json:"user"`
	Count int64  `bson:"count" json:"count"`
}

type ResponseIndividual struct {
	ResponseIndividual User
	StatusCode         int    `bson:"status_code" json:"status_code"`
	AccessToken        string `bson:"access_token" json:"access_token"`
}

//go:generate mockery --name IUserRepository
type IUserRepository interface {
	FetchMany(c context.Context) (Response, error)
	DeleteOne(c context.Context, userID string) error
	Login(c context.Context, request SignIn) (*User, error)
	GetByEmail(c context.Context, email string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
	Create(c context.Context, user User) error
	Update(ctx context.Context, userID string, user User) error
	UpsertOne(c context.Context, email string, user *User) (*User, error)
	UpdateImage(c context.Context, userID string, imageURL string) error
}
