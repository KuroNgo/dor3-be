package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

//go:generate mockery
type IUserRepository interface {
	Create(c context.Context, user *User) error
	CreateAsync(c context.Context, user *User) <-chan error
	Fetch(c context.Context) ([]User, error)
	Update(c context.Context, user *User) error
	Delete(userID primitive.ObjectID) error
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id primitive.ObjectID) (User, error)
}
