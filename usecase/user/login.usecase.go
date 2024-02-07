package usecase

import (
	user_domain "clean-architecture/domain/request/user"
	"context"
	"time"
)

type loginUseCase struct {
	userRepository user_domain.IUserRepository
	contextTimeout time.Duration
}

func NewLoginUseCase(userRepository user_domain.IUserRepository, timeout time.Duration) user_domain.ILoginUseCase {
	return &loginUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (l loginUseCase) GetAllUser(c context.Context) ([]user_domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (l loginUseCase) GetUserByEmail(c context.Context, email string) (user_domain.User, error) {
	ctx, cancel := context.WithTimeout(c, l.contextTimeout)
	defer cancel()
	return l.userRepository.GetByEmail(ctx, email)
}

func (l loginUseCase) GetUserByUsername(c context.Context, username string) (user_domain.User, error) {
	//TODO implement me
	panic("implement me")
}
