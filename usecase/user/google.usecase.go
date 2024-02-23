package user_usecase

import (
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal/oauth2/google"
	"context"
	"time"
)

type googleUseCase struct {
	userRepository user_domain.IUserRepository
	contextTimeout time.Duration
}

func NewGoogleUseCase(userRepository user_domain.IUserRepository, timeout time.Duration) user_domain.IGoogleAuthUseCase {
	return &googleUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (g *googleUseCase) UpsertUser(c context.Context, email string, user *user_domain.User) (*user_domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, g.contextTimeout)
	defer cancel()
	return g.userRepository.UpsertOne(ctx, email, user)
}

func (g *googleUseCase) GetGoogleOauthToken(code string) (*google.OauthToken, error) {
	return google.GetGoogleOauthToken(code)
}

func (g *googleUseCase) GetGoogleUser(accessToken string, idToken string) (*google.UserResult, error) {
	return google.GetGoogleUser(accessToken, idToken)
}
