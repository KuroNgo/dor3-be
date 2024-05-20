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

func (u *userAttemptUseCase) FetchOneByUnitIDAndUserID(ctx context.Context, userID string, unit string) (user_attempt_domain.UserProcess, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userAttemptUseCase) FetchOneByUnitID(ctx context.Context, unitID string) (user_attempt_domain.UserProcess, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userAttemptUseCase) CreateAttemptByExerciseID(ctx context.Context, userID user_attempt_domain.UserProcess) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.CreateAttemptByExerciseID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) UpdateAttemptByExamID(ctx context.Context, userID user_attempt_domain.UserProcess) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.UpdateAttemptByExamID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) UpdateAttemptByQuizID(ctx context.Context, userID user_attempt_domain.UserProcess) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.UpdateAttemptByQuizID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
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
