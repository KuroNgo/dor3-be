package unit_repo

import (
	unit_domain "clean-architecture/domain/_unit"
	"clean-architecture/infrastructor/mongo"
	"context"
)

type unitRepository struct {
	database         mongo.Database
	collectionUnit   string
	collectionLesson string
}

func (u *unitRepository) FetchByID(ctx context.Context, unitID string) (*unit_domain.Unit, error) {
	//TODO implement me
	panic("implement me")
}

func (u *unitRepository) FetchMany(ctx context.Context) ([]unit_domain.Unit, error) {
	//TODO implement me
	panic("implement me")
}

func (u *unitRepository) FetchToDeleteMany(ctx context.Context) (*[]unit_domain.Unit, error) {
	//TODO implement me
	panic("implement me")
}

func (u *unitRepository) CreateOne(ctx context.Context, unit *unit_domain.Unit) error {
	//TODO implement me
	panic("implement me")
}

func (u *unitRepository) UpdateOne(ctx context.Context, unitID string, unit unit_domain.Unit) error {
	//TODO implement me
	panic("implement me")
}

func (u *unitRepository) UpsertOne(ctx context.Context, id string, unit *unit_domain.Unit) (*unit_domain.Unit, error) {
	//TODO implement me
	panic("implement me")
}

func (u *unitRepository) DeleteOne(ctx context.Context, unitID string) error {
	//TODO implement me
	panic("implement me")
}

func NewUnitRepository(db mongo.Database, collectionUnit string, collectionLesson string) unit_domain.IUnitRepository {
	return &unitRepository{
		database:         db,
		collectionUnit:   collectionUnit,
		collectionLesson: collectionLesson,
	}
}
