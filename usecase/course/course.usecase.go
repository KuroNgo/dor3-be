package course_usecase

import (
	course_domain "clean-architecture/domain/course"
	"context"
	"time"
)

type courseUseCase struct {
	courseRepository course_domain.ICourseRepository
	contextTimeout   time.Duration
}

func NewCourseUseCase(courseRepository course_domain.ICourseRepository, timeout time.Duration) course_domain.ICourseUseCase {
	return &courseUseCase{
		courseRepository: courseRepository,
		contextTimeout:   timeout,
	}
}

func (c *courseUseCase) FetchByID(ctx context.Context, courseID string) (*course_domain.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, err := c.courseRepository.FetchByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return course, err
}

func (c *courseUseCase) FetchMany(ctx context.Context) ([]course_domain.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, err := c.courseRepository.FetchMany(ctx)
	if err != nil {
		return nil, err
	}

	return course, err
}

func (c *courseUseCase) FetchToDeleteMany(ctx context.Context) (*[]course_domain.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, err := c.courseRepository.FetchToDeleteMany(ctx)
	if err != nil {
		return nil, err
	}

	return course, err
}

func (c *courseUseCase) UpdateOne(ctx context.Context, courseID string, course course_domain.Course) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	err := c.courseRepository.UpdateOne(ctx, courseID, course)
	if err != nil {
		return err
	}

	return err
}

func (c *courseUseCase) CreateOne(ctx context.Context, course *course_domain.Course) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	err := c.courseRepository.CreateOne(ctx, course)

	if err != nil {
		return err
	}

	return nil
}

func (c *courseUseCase) UpsertOne(ctx context.Context, id string, course *course_domain.Course) (*course_domain.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	courseRes, err := c.courseRepository.UpsertOne(ctx, id, course)
	if err != nil {
		return nil, err
	}
	return courseRes, nil
}

func (c *courseUseCase) DeleteOne(ctx context.Context, courseID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	err := c.courseRepository.DeleteOne(ctx, courseID)
	if err != nil {
		return err
	}

	return err
}
