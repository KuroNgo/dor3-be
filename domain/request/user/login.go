package user_domain

import "context"

type LoginEmailRequest struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type LoginUsernameRequest struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}
type LoginResponse struct {
	AccessToken  string `json:"accessToken" bson:"access_token"`
	RefreshToken string `json:"refreshToken" bson:"refresh_token"`
}

type ILoginUseCase interface {
	GetCurrentUser(c context.Context) ([]User, error)
	GetUserByEmail(c context.Context, email string) (User, error)
	GetUserByUsername(c context.Context, username string) (User, error)
}
