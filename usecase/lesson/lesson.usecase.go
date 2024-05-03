package lesson_usecase

import (
	lesson_domain "clean-architecture/domain/lesson"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (l *lessonUseCase) FetchByID(ctx context.Context, lessonID string) (lesson_domain.LessonResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	data, err := l.lessonRepository.FetchByID(ctx, lessonID)

	if err != nil {
		return lesson_domain.LessonResponse{}, err
	}

	return data, nil
}

func (l *lessonUseCase) FindCourseIDByCourseName(ctx context.Context, courseName string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	data, err := l.lessonRepository.FindCourseIDByCourseName(ctx, courseName)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return data, nil
}

func (l *lessonUseCase) CreateOneByNameCourse(ctx context.Context, lesson *lesson_domain.Lesson) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	err := l.lessonRepository.CreateOneByNameCourse(ctx, lesson)

	if err != nil {
		return err
	}

	return nil
}

func (l *lessonUseCase) FetchByIdCourse(ctx context.Context, idCourse string) (lesson_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, err := l.lessonRepository.FetchByIdCourse(ctx, idCourse)
	if err != nil {
		return lesson_domain.Response{}, err
	}

	return lesson, err
}

func (l *lessonUseCase) UpdateComplete(ctx context.Context, lessonID string, lesson lesson_domain.Lesson) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	err := l.lessonRepository.UpdateComplete(ctx, lessonID, lesson)
	if err != nil {
		return err
	}

	return err
}

func (l *lessonUseCase) FetchMany(ctx context.Context) ([]lesson_domain.LessonResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, err := l.lessonRepository.FetchMany(ctx)
	if err != nil {
		return nil, err
	}

	return lesson, err
}

func (l *lessonUseCase) UpdateOne(ctx context.Context, lesson *lesson_domain.Lesson) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	data, err := l.lessonRepository.UpdateOne(ctx, lesson)
	if err != nil {
		return data, err
	}

	return data, err
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

func (l *lessonUseCase) DeleteOne(ctx context.Context, lessonID string) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	err := l.lessonRepository.DeleteOne(ctx, lessonID)
	if err != nil {
		return err
	}

	return err
}
