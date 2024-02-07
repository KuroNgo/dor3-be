package user_domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SignUpInput struct {
	Name     string `json:"name" bson:"name" binding:"required"`
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required,min=8"`
	Role     string `json:"role" bson:"role"`
	Provider string `json:"provider,omitempty" bson:"provider,omitempty"`
	Photo    string `json:"photo,omitempty" bson:"photo,omitempty"`
	Verified bool   `json:"verified" bson:"verified"`
}

type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type DBResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	PasswordConfirm string             `json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
	Provider        string             `json:"provider" bson:"provider"`
	Photo           string             `json:"photo,omitempty" bson:"photo,omitempty"`
	Role            string             `json:"role" bson:"role"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

type UpdateDBUser struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`
	Email           string             `json:"email,omitempty" bson:"email,omitempty"`
	Password        string             `json:"password,omitempty" bson:"password,omitempty"`
	PasswordConfirm string             `json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
	Role            string             `json:"role,omitempty" bson:"role,omitempty"`
	Provider        string             `json:"provider" bson:"provider"`
	Photo           string             `json:"photo,omitempty" bson:"photo,omitempty"`
	Verified        bool               `json:"verified,omitempty" bson:"verified,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func FilteredResponse(user *DBResponse) Response {
	return Response{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Provider:  user.Provider,
		Photo:     user.Photo,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
