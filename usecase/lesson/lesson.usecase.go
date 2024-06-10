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

func (l *lessonUseCase) FetchManyNotPaginationInUser(ctx context.Context, userID primitive.ObjectID) ([]lesson_domain.LessonProcessResponse, lesson_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, detail, err := l.lessonRepository.FetchManyNotPaginationInUser(ctx, userID)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}

	return lesson, detail, err
}

func (l *lessonUseCase) FetchByIDInUser(ctx context.Context, userID primitive.ObjectID, lessonID string) (lesson_domain.LessonProcessResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	data, err := l.lessonRepository.FetchByIDInUser(ctx, userID, lessonID)

	if err != nil {
		return lesson_domain.LessonProcessResponse{}, err
	}

	return data, nil
}

func (l *lessonUseCase) FetchByIDCourseInUser(ctx context.Context, userID primitive.ObjectID, courseID string, page string) ([]lesson_domain.LessonProcessResponse, lesson_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, detail, err := l.lessonRepository.FetchByIDCourseInUser(ctx, userID, courseID, page)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}

	return lesson, detail, err
}

func (l *lessonUseCase) FetchManyInUser(ctx context.Context, userID primitive.ObjectID, page string) ([]lesson_domain.LessonProcessResponse, lesson_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, detail, err := l.lessonRepository.FetchManyInUser(ctx, userID, page)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}

	return lesson, detail, err
}

func (l *lessonUseCase) FetchByIDInAdmin(ctx context.Context, lessonID string) (lesson_domain.LessonResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	data, err := l.lessonRepository.FetchByIDInAdmin(ctx, lessonID)

	if err != nil {
		return lesson_domain.LessonResponse{}, err
	}

	return data, nil
}

func (l *lessonUseCase) FetchManyNotPaginationInAdmin(ctx context.Context) ([]lesson_domain.LessonResponse, lesson_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	data, detail, err := l.lessonRepository.FetchManyNotPaginationInAdmin(ctx)

	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}

	return data, detail, nil
}

func (l *lessonUseCase) FetchByIdCourseInAdmin(ctx context.Context, idCourse string, page string) ([]lesson_domain.LessonResponse, lesson_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, detail, err := l.lessonRepository.FetchByIdCourseInAdmin(ctx, idCourse, page)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}

	return lesson, detail, err
}

func (l *lessonUseCase) FindLessonIDByLessonNameInAdmin(ctx context.Context, lessonName string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	data, err := l.lessonRepository.FindLessonIDByLessonNameInAdmin(ctx, lessonName)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return data, nil
}

func (l *lessonUseCase) FetchManyInAdmin(ctx context.Context, page string) ([]lesson_domain.LessonResponse, lesson_domain.DetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	lesson, detail, err := l.lessonRepository.FetchManyInAdmin(ctx, page)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}

	return lesson, detail, err
}

func (l *lessonUseCase) CreateOneByNameCourseInAdmin(ctx context.Context, lesson *lesson_domain.Lesson) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	err := l.lessonRepository.CreateOneByNameCourseInAdmin(ctx, lesson)

	if err != nil {
		return err
	}

	return nil
}

func (l *lessonUseCase) CreateOneInAdmin(ctx context.Context, lesson *lesson_domain.Lesson) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()
	err := l.lessonRepository.CreateOneInAdmin(ctx, lesson)

	if err != nil {
		return err
	}

	return nil
}

func (l *lessonUseCase) UpdateImageInAdmin(ctx context.Context, lesson *lesson_domain.Lesson) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	data, err := l.lessonRepository.UpdateImageInAdmin(ctx, lesson)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (l *lessonUseCase) UpdateOneInAdmin(ctx context.Context, lesson *lesson_domain.Lesson) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	data, err := l.lessonRepository.UpdateOneInAdmin(ctx, lesson)
	if err != nil {
		return data, err
	}

	return data, err
}

func (l *lessonUseCase) DeleteOneInAdmin(ctx context.Context, lessonID string) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	err := l.lessonRepository.DeleteOneInAdmin(ctx, lessonID)
	if err != nil {
		return err
	}

	return err
}
