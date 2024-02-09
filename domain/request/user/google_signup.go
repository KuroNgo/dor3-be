package user_domain

import (
	"clean-architecture/internal/Oauth2/google"
	"context"
)

type IGoogleAuthUseCase interface {
	UpsertUser(c context.Context, email string, user *User) (*Response, error)
	GetGoogleOauthToken(code string) (*google.OauthToken, error)
	GetGoogleUser(accessToken string, idToken string) (*google.UserResult, error)
}
