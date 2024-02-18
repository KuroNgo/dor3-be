package user_domain

import (
	"clean-architecture/internal/oauth2/google"
	"context"
)

//go:generate mockery --name IGoogleAuthUseCase
type IGoogleAuthUseCase interface {
	UpsertUser(c context.Context, email string, user *User) (*Response, error)
	GetGoogleOauthToken(code string) (*google.OauthToken, error)
	GetGoogleUser(accessToken string, idToken string) (*google.UserResult, error)
}
