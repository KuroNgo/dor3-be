package usecase

import (
	"clean-architecture/domain/request/user"
	"context"
	"time"
)

type loginUseCase struct {
	userRepository user.IUserRepository
	contextTimeout time.Duration
}

func NewLoginUseCase(userRepository user.IUserRepository, timeout time.Duration) user.ILoginUseCase {
	return &loginUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (lu *loginUseCase) GetUserByEmail(c context.Context, email string) (user.User, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByEmail(ctx, email)
}

func (lu *loginUseCase) GetUserByUsername(c context.Context, username string) (user.User, error) {
	//TODO implement me
	panic("implement me")
}
