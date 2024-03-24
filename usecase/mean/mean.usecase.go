package mean

import (
	mean_domain "clean-architecture/domain/mean"
	"context"
	"time"
)

type meanUseCase struct {
	meanRepository mean_domain.IMeanRepository
	contextTimeout time.Duration
}

func (m *meanUseCase) FetchMany(ctx context.Context) ([]mean_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	mean, err := m.meanRepository.FetchMany(ctx)
	if err != nil {
		return nil, err
	}

	return mean, err
}

func (m *meanUseCase) CreateOne(ctx context.Context, mean *mean_domain.Mean) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.meanRepository.CreateOne(ctx, mean)
	if err != nil {
		return err
	}

	return nil
}

func (m *meanUseCase) UpdateOne(ctx context.Context, meanID string, mean mean_domain.Mean) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.meanRepository.UpdateOne(ctx, meanID, mean)
	if err != nil {
		return err
	}

	return err
}

func (m *meanUseCase) UpsertOne(ctx context.Context, id string, mean *mean_domain.Mean) (*mean_domain.Mean, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	meanRes, err := m.meanRepository.UpsertOne(ctx, id, mean)
	if err != nil {
		return nil, err
	}
	return &meanRes, nil
}

func (m *meanUseCase) DeleteOne(ctx context.Context, meanID string) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.meanRepository.DeleteOne(ctx, meanID)
	if err != nil {
		return err
	}

	return err
}

func NewMeanUseCase(meanRepository mean_domain.IMeanRepository, timeout time.Duration) mean_domain.IMeanUseCase {
	return &meanUseCase{
		meanRepository: meanRepository,
		contextTimeout: timeout,
	}
}