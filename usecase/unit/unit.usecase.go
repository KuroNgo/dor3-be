package unit_usecase

import (
	unit_domain "clean-architecture/domain/unit"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (u *unitUseCase) FindLessonIDByLessonName(ctx context.Context, lessonName string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	data, err := u.unitRepository.FindLessonIDByLessonName(ctx, lessonName)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return data, err
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

func (u *unitUseCase) FetchByIdLesson(ctx context.Context, idLesson string) (unit_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	unit, err := u.unitRepository.FetchByIdLesson(ctx, idLesson)
	if err != nil {
		return unit_domain.Response{}, err
	}

	return unit, err
}

func (u *unitUseCase) UpdateComplete(ctx context.Context, update unit_domain.Update) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	err := u.unitRepository.UpdateComplete(ctx, update)
	if err != nil {
		return err
	}

	return nil
}

func (u *unitUseCase) FetchMany(ctx context.Context) (unit_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	unit, err := u.unitRepository.FetchMany(ctx)
	if err != nil {
		return unit_domain.Response{}, err
	}

	return unit, err
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

func (u *unitUseCase) UpdateOne(ctx context.Context, unitID string, unit unit_domain.Unit) (unit_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	unitRes, err := u.unitRepository.UpsertOne(ctx, unitID, &unit)
	if err != nil {
		return unit_domain.Response{}, err
	}
	return unitRes, nil
}

func (u *unitUseCase) UpsertOne(ctx context.Context, id string, unit *unit_domain.Unit) (unit_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	unitRes, err := u.unitRepository.UpsertOne(ctx, id, unit)
	if err != nil {
		return unit_domain.Response{}, err
	}
	return unitRes, nil
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
