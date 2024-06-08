package unit_usecase

import (
	unit_domain "clean-architecture/domain/unit"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type unitUseCase struct {
	unitRepository unit_domain.IUnitRepository
	contextTimeout time.Duration
}

func NewUnitUseCase(unitRepository unit_domain.IUnitRepository, timeout time.Duration) unit_domain.IUnitUseCase {
	return &unitUseCase{
		unitRepository: unitRepository,
		contextTimeout: timeout,
	}
}

func (u *unitUseCase) FetchMany(ctx context.Context, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	unit, detail, err := u.unitRepository.FetchMany(ctx, page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	return unit, detail, err
}

func (u *unitUseCase) FetchManyNotPagination(ctx context.Context) ([]unit_domain.UnitResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	unit, err := u.unitRepository.FetchManyNotPagination(ctx)
	if err != nil {
		return nil, err
	}

	return unit, err
}

func (u *unitUseCase) FetchByIdLesson(ctx context.Context, idLesson string, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	unit, detail, err := u.unitRepository.FetchByIdLesson(ctx, idLesson, page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	return unit, detail, err
}

func (u *unitUseCase) FetchOneByID(ctx context.Context, id string) (unit_domain.UnitResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	data, err := u.unitRepository.FetchOneByID(ctx, id)

	if err != nil {
		return unit_domain.UnitResponse{}, err
	}

	return data, nil
}

func (u *unitUseCase) CreateOne(ctx context.Context, unit *unit_domain.Unit) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.unitRepository.CreateOne(ctx, unit)

	if err != nil {
		return err
	}

	return nil
}

func (u *unitUseCase) CreateOneByNameLesson(ctx context.Context, unit *unit_domain.Unit) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.unitRepository.CreateOneByNameLesson(ctx, unit)

	if err != nil {
		return err
	}

	return nil
}

func (u *unitUseCase) UpdateOne(ctx context.Context, unit *unit_domain.Unit) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	data, err := u.unitRepository.UpdateOne(ctx, unit)

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *unitUseCase) DeleteOne(ctx context.Context, unitID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.unitRepository.DeleteOne(ctx, unitID)
	if err != nil {
		return err
	}

	return err
}
