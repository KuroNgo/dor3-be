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
	Username   string             `bson:"username"  json:"username"`
	Password   string             `bson:"password"  json:"password"`
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
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	Photo     string             `json:"photo,omitempty" bson:"photo,omitempty"`
	Provider  string             `json:"provider" bson:"provider"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

//go:generate mockery --name IUserRepository
type IUserRepository interface {
	Fetch(c context.Context) ([]User, error)
	Update(c context.Context, userID primitive.ObjectID, updatedUser interface{}) error
	Delete(c context.Context, userID primitive.ObjectID) error
	GetByEmail(c context.Context, email string) (User, error)
	GetByUsername(c context.Context, username string) (User, error)
	GetByID(c context.Context, id primitive.ObjectID) (User, error)
	UpsertUser(c context.Context, email string, user *User) (*Response, error)
}
