package user_usecase

import (
	user_domain "clean-architecture/domain/user"
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

func (g *googleUseCase) UpsertUser(c context.Context, email string, user *user_domain.UserInput) (*user_domain.User, error) {
	ctx, cancel := context.WithTimeout(c, g.contextTimeout)
	defer cancel()
	return g.userRepository.UpsertOne(ctx, email, user)
}
