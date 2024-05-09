package user_attempt_usecase

import (
	user_attempt_domain "clean-architecture/domain/user_attempt"
	"context"
	"time"
)

type userAttemptUseCase struct {
	userAttemptRepository user_attempt_domain.IUserProcessRepository
	contextTimeout        time.Duration
}

func NewAttemptUseCase(userAttemptRepository user_attempt_domain.IUserProcessRepository, timeout time.Duration) user_attempt_domain.IUserProcessUseCase {
	return &userAttemptUseCase{
		userAttemptRepository: userAttemptRepository,
		contextTimeout:        timeout,
	}
}

func (u *userAttemptUseCase) FetchManyByUserID(ctx context.Context, userID string) (user_attempt_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userAttempt, err := u.userAttemptRepository.FetchManyByUserID(ctx, userID)
	if err != nil {
		return user_attempt_domain.Response{}, err
	}

	return userAttempt, err
}

func (u *userAttemptUseCase) CreateOneByUserID(ctx context.Context, user user_attempt_domain.UserProcess) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.CreateOneByUserID(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) UpdateAttemptByUserID(ctx context.Context, user user_attempt_domain.UserProcess) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.UpdateAttemptByUserID(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) DeleteAllByUserID(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.DeleteAllByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
