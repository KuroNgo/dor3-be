package mark_list_usecase

import (
	markList_domain "clean-architecture/domain/mark_list"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type markListUseCase struct {
	markListRepository markList_domain.IMarkListRepository
	contextTimeout     time.Duration
}

func NewMarkListUseCase(markListRepository markList_domain.IMarkListRepository, timeout time.Duration) markList_domain.IMarkListUseCase {
	return &markListUseCase{
		markListRepository: markListRepository,
		contextTimeout:     timeout,
	}
}

func (m *markListUseCase) FetchManyByUser(ctx context.Context, user primitive.ObjectID) (markList_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	markList, err := m.markListRepository.FetchManyByUser(ctx, user)
	if err != nil {
		return markList_domain.Response{}, err
	}

	return markList, err
}

func (m *markListUseCase) FetchByIdByUser(ctx context.Context, user primitive.ObjectID, id string) (markList_domain.MarkList, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	markList, err := m.markListRepository.FetchByIdByUser(ctx, user, id)
	if err != nil {
		return markList_domain.MarkList{}, err
	}

	return markList, err
}

func (m *markListUseCase) FetchManyByUserIDByUser(ctx context.Context, user primitive.ObjectID) (markList_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	markList, err := m.markListRepository.FetchManyByUser(ctx, user)
	if err != nil {
		return markList_domain.Response{}, err
	}

	return markList, err
}

func (m *markListUseCase) UpdateOneByUser(ctx context.Context, user primitive.ObjectID, markList *markList_domain.MarkList) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	data, err := m.markListRepository.UpdateOneByUser(ctx, user, markList)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *markListUseCase) CreateOneByUser(ctx context.Context, user primitive.ObjectID, markList *markList_domain.MarkList) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.markListRepository.CreateOneByUser(ctx, user, markList)
	if err != nil {
		return err
	}

	return nil
}

func (m *markListUseCase) DeleteOneByUser(ctx context.Context, user primitive.ObjectID, markListID string) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.markListRepository.DeleteOneByUser(ctx, user, markListID)
	if err != nil {
		return err
	}

	return nil
}

func (m *markListUseCase) FetchManyByAdmin(ctx context.Context) (markList_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	markList, err := m.markListRepository.FetchManyByAdmin(ctx)
	if err != nil {
		return markList_domain.Response{}, err
	}

	return markList, err
}

func (m *markListUseCase) FetchByIdByAdmin(ctx context.Context, id string) (markList_domain.MarkList, error) {
	//TODO implement me
	panic("implement me")
}
