package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

type SignInWithUsername struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type SignUpInput struct {
	FullName string `bson:"fullName" json:"fullName"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
	Phone    string `bson:"phone" json:"phone"`
}

type SignInWithEmail struct {
	Email    string `json:"email"  binding:"required" bson:"email"`
	Password string `json:"password"  binding:"required" bson:"password"`
}
