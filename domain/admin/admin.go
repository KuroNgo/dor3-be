package admin_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionAdmin = "admin"
)

type Admin struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	FullName  string             `bson:"full_name" json:"full_name"`
	Password  string             `bson:"password" json:"password"`
	Avatar    string             `bson:"avatar" json:"avatar"`
	Address   string             `bson:"address" json:"address"`
	Role      string             `bson:"role" json:"role"`
	Phone     string             `bson:"phone" json:"phone"`
	Email     string             `bson:"email" json:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Response struct {
	Admin []Admin
}

type IAdminRepository interface {
	Login(c context.Context, request SignIn) (*Admin, error)
	GetByID(c context.Context, id string) (*Admin, error)
	FetchMany(ctx context.Context) (Response, error)
	GetByEmail(ctx context.Context, username string) (*Admin, error)
	CreateOne(ctx context.Context, admin Admin) error
	UpdateOne(ctx context.Context, admin *Admin) (*mongo.UpdateResult, error)
	ChangeEmail(ctx context.Context, email string) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, adminID string) error
	UpsertOne(ctx context.Context, email string, admin *Admin) error
}
