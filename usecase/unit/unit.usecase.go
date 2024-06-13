package unit_usecase

import (
	unit_domain "clean-architecture/domain/unit"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type unitUseCase struct {
	unitRepository unit_domain.IUnitRepository
	contextTimeout time.Duration
}

func (u *unitUseCase) FetchManyInUser(ctx context.Context, user primitive.ObjectID, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *unitUseCase) FetchOneByIDInUser(ctx context.Context, user primitive.ObjectID, id string) (unit_domain.UnitResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *unitUseCase) FetchManyNotPaginationInUser(ctx context.Context, user primitive.ObjectID) ([]unit_domain.UnitResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *unitUseCase) FetchByIdLessonInUser(ctx context.Context, user primitive.ObjectID, idLesson string, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewUnitUseCase(unitRepository unit_domain.IUnitRepository, timeout time.Duration) unit_domain.IUnitUseCase {
	return &unitUseCase{
		unitRepository: unitRepository,
		contextTimeout: timeout,
	}
}

func (u *unitUseCase) FetchManyInAdmin(ctx context.Context, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	unit, detail, err := u.unitRepository.FetchManyInAdmin(ctx, page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	return unit, detail, err
}

func (u *unitUseCase) FetchManyNotPaginationInAdmin(ctx context.Context) ([]unit_domain.UnitResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	unit, err := u.unitRepository.FetchManyNotPaginationInAdmin(ctx)
	if err != nil {
		return nil, err
	}

	return unit, err
}

func (u *unitUseCase) FetchByIdLessonInAdmin(ctx context.Context, idLesson string, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	unit, detail, err := u.unitRepository.FetchByIdLessonInAdmin(ctx, idLesson, page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	return unit, detail, err
}

func (u *unitUseCase) FetchOneByIDInAdmin(ctx context.Context, id string) (unit_domain.UnitResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	data, err := u.unitRepository.FetchOneByIDInAdmin(ctx, id)

	if err != nil {
		return unit_domain.UnitResponse{}, err
	}

	return data, nil
}

func (u *unitUseCase) CreateOneInAdmin(ctx context.Context, unit *unit_domain.Unit) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.unitRepository.CreateOneInAdmin(ctx, unit)

	if err != nil {
		return err
	}

	return nil
}

func (u *unitUseCase) CreateOneByNameLessonInAdmin(ctx context.Context, unit *unit_domain.Unit) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.unitRepository.CreateOneByNameLessonInAdmin(ctx, unit)

	if err != nil {
		return err
	}

	return nil
}

func (u *unitUseCase) UpdateOneInAdmin(ctx context.Context, unit *unit_domain.Unit) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	data, err := u.unitRepository.UpdateOneInAdmin(ctx, unit)

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *unitUseCase) UpdateCompleteInUser(ctx context.Context, user primitive.ObjectID) (*mongo.UpdateResult, error) {
	//TODO implement me
	panic("implement me")
}

func (u *unitUseCase) DeleteOneInAdmin(ctx context.Context, unitID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.unitRepository.DeleteOneInAdmin(ctx, unitID)
	if err != nil {
		return err
	}

	return err
}
