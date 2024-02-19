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
	ID         primitive.ObjectID `bson:"_id" json:"-"`
	FullName   string             `bson:"full_name"  json:"title"`
	Email      string             `bson:"email"  json:"email"`
	AvatarURL  string             `bson:"avatar_url"  json:"avatar_url"`
	Specialize string             `bson:"specialize"  json:"specialize"`
	Phone      string             `bson:"phone"   json:"phone"`
	Age        uint8              `bson:"age"  json:"age"`
	Provider   string             `json:"provider,omitempty" bson:"provider,omitempty"`
	Verified   bool               `json:"verified" bson:"verified"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	Role       string             `bson:"role" json:"role"`
}

type Response struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Role      string             `json:"role" bson:"role"`
	Photo     string             `json:"photo" bson:"photo"`
	Provider  string             `json:"provider" bson:"provider"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

//go:generate mockery --name IUserRepository
type IUserRepository interface {
	Fetch(c context.Context) ([]User, error)
	Delete(c context.Context, userID string) error
	GetByEmail(c context.Context, email string) (*User, error)
	GetByUsername(c context.Context, username string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
	Upsert(c context.Context, email string, user *User) (*Response, error)
}
