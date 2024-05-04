package course_usecase

import (
	course_domain "clean-architecture/domain/course"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
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

func (c *courseUseCase) FetchByID(ctx context.Context, courseID string) (course_domain.CourseResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, err := c.courseRepository.FetchByID(ctx, courseID)
	if err != nil {
		return course_domain.CourseResponse{}, err
	}

	return course, err
}

func (c *courseUseCase) FetchManyForEachCourse(ctx context.Context, page string) ([]course_domain.CourseResponse, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, count, err := c.courseRepository.FetchManyForEachCourse(ctx, page)
	if err != nil {
		return []course_domain.CourseResponse{}, 0, err
	}

	return course, count, err
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

func (c *courseUseCase) UpdateOne(ctx context.Context, course *course_domain.Course) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	data, err := c.courseRepository.UpdateOne(ctx, course)
	if err != nil {
		return nil, err
	}

	return data, nil
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
