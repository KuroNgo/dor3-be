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
	ID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	FullName   string             `bson:"full_name"  json:"full_name"`
	Email      string             `bson:"email"  json:"email"`
	AvatarURL  string             `bson:"avatar_url"  json:"avatar_url"`
	Specialize string             `bson:"specialize"  json:"specialize"`
	Phone      string             `bson:"phone"   json:"phone"`
	Age        uint8              `bson:"age"  json:"age"`
	Provider   string             `json:"provider" bson:"provider"`
	Verified   bool               `json:"verified" bson:"verified"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	Role       string             `bson:"role" json:"role"`
}

type Response struct {
	ID         primitive.ObjectID `json:"id" bson:"id"`
	FullName   string             `json:"full_name"  bson:"full_name"`
	Email      string             `json:"email" bson:"email"`
	AvatarURL  string             `json:"avatar_url"  bson:"avatar_url"`
	Specialize string             `json:"specialize"  bson:"specialize"`
	Role       string             `json:"role" bson:"role"`
	Photo      string             `json:"photo" bson:"photo"`
	Provider   string             `json:"provider" bson:"provider"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

//go:generate mockery --name IUserRepository
type IUserRepository interface {
	FetchMany(c context.Context) ([]User, error)
	DeleteOne(c context.Context, userID string) error
	GetByEmail(c context.Context, email string) (*User, error)
	GetByUsername(c context.Context, username string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
	Create(c context.Context, user Input) error
	UpsertOne(c context.Context, email string, user *User) (*Response, error)
}
