package user_attempt_usecase

import (
	"clean-architecture/domain/user_process"
	"context"
	"time"
)

type userAttemptUseCase struct {
	userAttemptRepository user_process.IUserProcessRepository
	contextTimeout        time.Duration
}

func NewAttemptUseCase(userAttemptRepository user_process.IUserProcessRepository, timeout time.Duration) user_process.IUserProcessUseCase {
	return &userAttemptUseCase{
		userAttemptRepository: userAttemptRepository,
		contextTimeout:        timeout,
	}
}

func (u *userAttemptUseCase) FetchOneByUnitIDAndUserID(ctx context.Context, userID string, unit string) (user_process.ExamManagement, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userAttemptUseCase) FetchOneByUnitID(ctx context.Context, unitID string) (user_process.ExamManagement, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userAttemptUseCase) CreateExamManagementByExerciseID(ctx context.Context, userID user_process.ExamManagement) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.CreateExamManagementByExerciseID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) UpdateExamManagementByExamID(ctx context.Context, userID user_process.ExamManagement) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.UpdateExamManagementByExamID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) UpdateExamManagementByQuizID(ctx context.Context, userID user_process.ExamManagement) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.UpdateExamManagementByQuizID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) FetchManyByUserID(ctx context.Context, userID string) (user_process.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userAttempt, err := u.userAttemptRepository.FetchManyByUserID(ctx, userID)
	if err != nil {
		return user_process.Response{}, err
	}

	return userAttempt, err
}

func (u *userAttemptUseCase) UpdateExamManagementByUserID(ctx context.Context, user user_process.ExamManagement) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.UpdateExamManagementByUserID(ctx, user)
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
