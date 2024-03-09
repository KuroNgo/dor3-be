package lesson_usecase

import (
	lesson_domain "clean-architecture/domain/lesson"
	"context"
	"time"
)

type lessonUseCase struct {
	lessonRepository lesson_domain.ILessonRepository
	contextTimeout   time.Duration
}

func NewLessonUseCase(lessonRepository lesson_domain.ILessonRepository, timeout time.Duration) lesson_domain.ILessonUseCase {
	return &lessonUseCase{
		lessonRepository: lessonRepository,
		contextTimeout:   timeout,
	}
}

func (l *lessonUseCase) FetchByID(ctx context.Context, lessonID string) (*lesson_domain.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	course, err := l.lessonRepository.FetchByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	return course, err
}

func (l *lessonUseCase) FetchMany(ctx context.Context) ([]lesson_domain.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, err := l.lessonRepository.FetchMany(ctx)
	if err != nil {
		return nil, err
	}

	return lesson, err
}

func (l *lessonUseCase) FetchToDeleteMany(ctx context.Context) (*[]lesson_domain.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, err := l.lessonRepository.FetchToDeleteMany(ctx)
	if err != nil {
		return nil, err
	}

	return lesson, err
}

func (l *lessonUseCase) UpdateOne(ctx context.Context, lessonID string, lesson lesson_domain.Lesson) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	err := l.lessonRepository.UpdateOne(ctx, lessonID, lesson)
	if err != nil {
		return err
	}

	return err
}

func (l *lessonUseCase) CreateOne(ctx context.Context, lesson *lesson_domain.Lesson) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	err := l.lessonRepository.CreateOne(ctx, lesson)

	if err != nil {
		return err
	}

	return nil
}

func (l *lessonUseCase) UpsertOne(ctx context.Context, id string, lesson *lesson_domain.Lesson) (*lesson_domain.Lesson, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, err := l.lessonRepository.UpsertOne(ctx, id, lesson)
	if err != nil {
		return nil, err
	}
	return lesson, nil
}

func (l *lessonUseCase) DeleteOne(ctx context.Context, lessonID string) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	err := l.lessonRepository.DeleteOne(ctx, lessonID)
	if err != nil {
		return err
	}

	return err
}
