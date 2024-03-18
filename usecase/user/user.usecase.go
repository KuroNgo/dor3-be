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

func (u *userUseCase) Create(c context.Context, user user_domain.User) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	err := u.userRepository.Create(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) Login(c context.Context, request user_domain.SignIn) (*user_domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *userUseCase) Fetch(c context.Context) ([]user_domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.FetchMany(ctx)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *userUseCase) GetByEmail(c context.Context, email string) (*user_domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, err
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

func (u *userUseCase) Update(ctx context.Context, userID string, user user_domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepository.Update(ctx, userID, user)

	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) Delete(ctx context.Context, userID string, user user_domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepository.Update(ctx, userID, user)

	if err != nil {
		return err
	}

	return nil
}
