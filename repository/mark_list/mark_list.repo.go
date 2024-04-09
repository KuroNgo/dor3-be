package mark_list

import (
	mark_list_repository "clean-architecture/domain/mark_list"
	"clean-architecture/infrastructor/mongo"
	"context"
)

type markListRepository struct {
	database           mongo.Database
	collectionMarkList string
}

func (m markListRepository) FetchMany(ctx context.Context) (mark_list_repository.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (m markListRepository) UpdateOne(ctx context.Context, markListID string, markList mark_list_repository.MarkList) error {
	//TODO implement me
	panic("implement me")
}

func (m markListRepository) CreateOne(ctx context.Context, markList *mark_list_repository.MarkList) error {
	//TODO implement me
	panic("implement me")
}

func (m markListRepository) UpsertOne(c context.Context, id string, markList *mark_list_repository.MarkList) (mark_list_repository.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (m markListRepository) DeleteOne(ctx context.Context, markListID string) error {
	//TODO implement me
	panic("implement me")
}

func NewExerciseRepository(db mongo.Database, collectionMarkList string) mark_list_repository.IMarkListRepository {
	return &markListRepository{
		database:           db,
		collectionMarkList: collectionMarkList,
	}
}
