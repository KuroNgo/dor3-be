package admin_usecase

import (
	admin_domain "clean-architecture/domain/admin"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type adminUseCase struct {
	adminRepository admin_domain.IAdminRepository
	contextTimeout  time.Duration
}

func NewAdminUseCase(adminRepository admin_domain.IAdminRepository, timeout time.Duration) admin_domain.IAdminUseCase {
	return &adminUseCase{
		adminRepository: adminRepository,
		contextTimeout:  timeout,
	}
}

func (a *adminUseCase) GetByID(ctx context.Context, id string) (*admin_domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	quiz, err := a.adminRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return quiz, err
}

func (a *adminUseCase) Login(c context.Context, request admin_domain.SignIn) (*admin_domain.Admin, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	quiz, err := a.adminRepository.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return quiz, err
}

func (a *adminUseCase) FetchMany(ctx context.Context) (admin_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	quiz, err := a.adminRepository.FetchMany(ctx)
	if err != nil {
		return admin_domain.Response{}, err
	}

	return quiz, err
}

func (a *adminUseCase) GetByEmail(ctx context.Context, username string) (*admin_domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	admin, err := a.adminRepository.GetByEmail(ctx, username)
	if err != nil {
		return nil, err
	}

	return admin, err
}

func (a *adminUseCase) CreateOne(ctx context.Context, admin admin_domain.Admin) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.adminRepository.CreateOne(ctx, admin)
	if err != nil {
		return err
	}

	return err
}

func (a *adminUseCase) UpdateOne(ctx context.Context, admin *admin_domain.Admin) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	data, err := a.adminRepository.UpdateOne(ctx, admin)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (a *adminUseCase) DeleteOne(ctx context.Context, adminID string) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.adminRepository.DeleteOne(ctx, adminID)
	if err != nil {
		return err
	}

	return err
}

func (a *adminUseCase) UpsertOne(ctx context.Context, email string, admin *admin_domain.Admin) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.adminRepository.UpsertOne(ctx, email, admin)
	if err != nil {
		return err
	}

	return err
}
