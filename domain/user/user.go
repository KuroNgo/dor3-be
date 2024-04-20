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
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	FullName  string             `bson:"full_name"  json:"full_name"`
	Email     string             `bson:"email"  json:"email"`
	Password  string             `bson:"password"  json:"password"`
	AvatarURL string             `bson:"avatar_url"  json:"avatar_url"`
	AssetID   string             `bson:"asset_id"  json:"asset_id"`
	Phone     string             `bson:"phone"   json:"phone"`
	Provider  string             `bson:"provider" json:"provider"`
	Verified  bool               `bson:"verified" json:"verified"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Role      string             `bson:"role" json:"role"`
}

type Response struct {
	Count int64  `bson:"count" json:"count"`
	User  []User `bson:"user" json:"user"`
}

//go:generate mockery --name IUserRepository
type IUserRepository interface {
	FetchMany(c context.Context) (Response, error)
	DeleteOne(c context.Context, userID string) error
	Login(c context.Context, request SignIn) (*User, error)
	GetByEmail(c context.Context, email string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
	Create(c context.Context, user User) error
	Update(ctx context.Context, user *User) (*mongo.UpdateResult, error)
	UpsertOne(c context.Context, email string, user *User) (*User, error)
	UpdateImage(c context.Context, userID string, imageURL string) error
}
