package mark_list_usecase

import (
	markList_domain "clean-architecture/domain/mark_list"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type markListUseCase struct {
	markListRepository markList_domain.IMarkListRepository
	contextTimeout     time.Duration
}

func (m *markListUseCase) FetchManyByUserID(ctx context.Context, userId string) (markList_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	markList, err := m.markListRepository.FetchManyByUserID(ctx, userId)
	if err != nil {
		return markList_domain.Response{}, err
	}

	return markList, err
}

func (m *markListUseCase) UpdateOne(ctx context.Context, markList *markList_domain.MarkList) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	data, err := m.markListRepository.UpdateOne(ctx, markList)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *markListUseCase) CreateOne(ctx context.Context, markList *markList_domain.MarkList) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.markListRepository.CreateOne(ctx, markList)
	if err != nil {
		return err
	}

	return nil
}

func (m *markListUseCase) DeleteOne(ctx context.Context, markListID string) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.markListRepository.DeleteOne(ctx, markListID)
	if err != nil {
		return err
	}

	return nil
}

func NewMarkListUseCase(markListRepository markList_domain.IMarkListRepository, timeout time.Duration) markList_domain.IMarkListUseCase {
	return &markListUseCase{
		markListRepository: markListRepository,
		contextTimeout:     timeout,
	}
}
