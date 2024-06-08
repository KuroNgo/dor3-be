package user_attempt_usecase

import (
	"clean-architecture/domain/user_process/exam_management"
	"context"
	"time"
)

type userAttemptUseCase struct {
	userAttemptRepository exam_management.IUserProcessRepository
	contextTimeout        time.Duration
}

func NewAttemptUseCase(userAttemptRepository exam_management.IUserProcessRepository, timeout time.Duration) exam_management.IUserProcessUseCase {
	return &userAttemptUseCase{
		userAttemptRepository: userAttemptRepository,
		contextTimeout:        timeout,
	}
}

func (u *userAttemptUseCase) FetchOneByUnitIDAndUserID(ctx context.Context, userID string, unit string) (exam_management.ExamManagement, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userAttemptUseCase) FetchOneByUnitID(ctx context.Context, unitID string) (exam_management.ExamManagement, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userAttemptUseCase) CreateExamManagementByExerciseID(ctx context.Context, userID exam_management.ExamManagement) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.CreateExamManagementByExerciseID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) UpdateExamManagementByExamID(ctx context.Context, userID exam_management.ExamManagement) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.UpdateExamManagementByExamID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) UpdateExamManagementByQuizID(ctx context.Context, userID exam_management.ExamManagement) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.userAttemptRepository.UpdateExamManagementByQuizID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userAttemptUseCase) FetchManyByUserID(ctx context.Context, userID string) (exam_management.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userAttempt, err := u.userAttemptRepository.FetchManyByUserID(ctx, userID)
	if err != nil {
		return exam_management.Response{}, err
	}

	return userAttempt, err
}

func (u *userAttemptUseCase) UpdateExamManagementByUserID(ctx context.Context, user exam_management.ExamManagement) error {
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
