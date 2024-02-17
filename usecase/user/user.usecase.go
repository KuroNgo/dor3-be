package user_usecase

import (
	user_domain "clean-architecture/domain/user"
	"context"
	"time"
)

type userUseCase struct {
	userRepository user_domain.IUserRepository
	contextTimeout time.Duration
}

func NewUserUseCase(userRepository user_domain.IUserRepository, timeout time.Duration) user_domain.IUserUseCase {
	return &userUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (u *userUseCase) Fetch(c context.Context) ([]user_domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepository.Fetch(ctx)
}

func (u *userUseCase) GetByEmail(c context.Context, email string) (*user_domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepository.GetByEmail(ctx, email)
}

func (u *userUseCase) GetByID(c context.Context, id string) (*user_domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, err
}