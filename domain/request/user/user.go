package user

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
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	RoleID     primitive.ObjectID `bson:"_id" json:"-"`
}

type IUserRepository interface {
	Create(c context.Context, user *User) error
	CreateAsync(c context.Context, user *User) <-chan error
	Fetch(c context.Context) ([]User, error)
	Update(c context.Context, userID primitive.ObjectID, updatedUser interface{}) error
	Delete(c context.Context, userID primitive.ObjectID) error
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id primitive.ObjectID) (User, error)
}
