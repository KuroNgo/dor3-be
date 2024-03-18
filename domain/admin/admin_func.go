package admin_domain

import "context"

type SignUp struct {
	Username string `bson:"username" json:"username"`
	FullName string `bson:"full_name" json:"full_name"`
	Password string `bson:"password" json:"password"`
	Avatar   string `bson:"avatar" json:"avatar"`
	Address  string `bson:"address" json:"address"`
	Phone    string `bson:"phone" json:"phone"`
	Email    string `bson:"email" json:"email"`
}

type SignIn struct {
	Email    string `json:"email" bson:"email"`
	Password string `bson:"password"  json:"password"`
}

type IAdminUseCase interface {
	Login(c context.Context, request SignIn) (*Admin, error)
	GetByID(ctx context.Context, id string) (*Admin, error)
	FetchMany(ctx context.Context) ([]Admin, error)
	GetByEmail(ctx context.Context, email string) (*Admin, error)
	CreateOne(ctx context.Context, admin Admin) error
	UpdateOne(ctx context.Context, adminID string, admin Admin) error
	DeleteOne(ctx context.Context, adminID string) error
	UpsertOne(ctx context.Context, email string, admin *Admin) error
}
