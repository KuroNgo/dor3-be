package audio

import (
	audio_domain "clean-architecture/domain/audio"
	"context"
	"time"
)

type audioUseCase struct {
	audioRepository audio_domain.IAudioRepository
	contextTimeout  time.Duration
}

func NewAudioUseCase(audioRepository audio_domain.IAudioUseCase, timeout time.Duration) audio_domain.IAudioUseCase {
	return &audioUseCase{
		audioRepository: audioRepository,
		contextTimeout:  timeout,
	}
}

func (a *audioUseCase) FetchByID(ctx context.Context, audioID string) (*audio_domain.Audio, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	audio, err := a.audioRepository.FetchByID(ctx, audioID)
	if err != nil {
		return nil, err
	}

	return audio, err
}

func (a *audioUseCase) FetchMany(ctx context.Context) ([]audio_domain.Audio, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	quiz, err := a.audioRepository.FetchMany(ctx)
	if err != nil {
		return nil, err
	}

	return quiz, err
}

func (a *audioUseCase) FetchToDeleteMany(ctx context.Context) (*[]audio_domain.Audio, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	user, err := a.audioRepository.FetchToDeleteMany(ctx)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (a *audioUseCase) UpdateOne(ctx context.Context, audioID string, audio audio_domain.Audio) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.audioRepository.UpdateOne(ctx, audioID, audio)
	if err != nil {
		return err
	}

	return err
}

func (a *audioUseCase) CreateOne(ctx context.Context, audio *audio_domain.AutoMatch) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	err := a.audioRepository.CreateOne(ctx, audio)

	if err != nil {
		return err
	}

	return nil
}

func (a *audioUseCase) DeleteOne(ctx context.Context, audioID string) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.audioRepository.DeleteOne(ctx, audioID)
	if err != nil {
		return err
	}

	return err
}

func (a *audioUseCase) DeleteMany(ctx context.Context, audioIDs ...string) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	err := a.audioRepository.DeleteMany(ctx, audioIDs...)
	if err != nil {
		return err
	}

	return err
}
