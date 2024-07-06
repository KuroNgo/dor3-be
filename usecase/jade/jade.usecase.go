package jade_usecase

import (
	jade_domain "clean-architecture/domain/jade"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type jadeUseCase struct {
	jadeRepository jade_domain.IJadeRepository
	contextTimeout time.Duration
}

func NewJadeUseCase(jadeRepository jade_domain.IJadeRepository, timeout time.Duration) jade_domain.IJadeUseCase {
	return &jadeUseCase{
		jadeRepository: jadeRepository,
		contextTimeout: timeout,
	}
}

func (j *jadeUseCase) FetchJadeInUser(ctx context.Context, userID primitive.ObjectID) (jade_domain.JadeBlockchain, error) {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	jade, err := j.jadeRepository.FetchJadeInUser(ctx, userID)
	if err != nil {
		return jade_domain.JadeBlockchain{}, err
	}

	return jade, err
}

func (j *jadeUseCase) CreateCurrencyInUser(ctx context.Context, userID primitive.ObjectID, data int64) error {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	err := j.jadeRepository.CreateCurrencyInUser(ctx, userID, data)
	if err != nil {
		return err
	}

	return err
}

func (j *jadeUseCase) UpdateCurrencyInUser(ctx context.Context, userID primitive.ObjectID, data int64) error {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	err := j.jadeRepository.UpdateCurrencyInUser(ctx, userID, data)
	if err != nil {
		return err
	}

	return err
}

func (j *jadeUseCase) Rank(ctx context.Context) ([]jade_domain.JadeBlockchain, error) {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	jade, err := j.jadeRepository.Rank(ctx)
	if err != nil {
		return nil, err
	}

	return jade, err
}
