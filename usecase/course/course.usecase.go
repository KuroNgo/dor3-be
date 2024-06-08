package course_usecase

import (
	course_domain "clean-architecture/domain/course"
	lesson_management_domain "clean-architecture/domain/user_process/lesson_management"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (c *courseUseCase) UpdateCompleteInUser(ctx context.Context) (*mongo.UpdateResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *courseUseCase) FetchManyInUser(ctx context.Context, userID primitive.ObjectID, page string) ([]lesson_management_domain.CourseProcess, course_domain.DetailForManyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, detail, err := c.courseRepository.FetchManyInUser(ctx, userID, page)
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}

	return course, detail, nil
}

func (c *courseUseCase) FetchByIDInUser(ctx context.Context, userID primitive.ObjectID, courseID string) (lesson_management_domain.CourseProcess, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, err := c.courseRepository.FetchByIDInUser(ctx, userID, courseID)
	if err != nil {
		return lesson_management_domain.CourseProcess{}, err
	}

	return course, nil
}

func (c *courseUseCase) FetchByIDInAdmin(ctx context.Context, courseID string) (course_domain.CourseResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, err := c.courseRepository.FetchByIDInAdmin(ctx, courseID)
	if err != nil {
		return course_domain.CourseResponse{}, err
	}

	return course, err
}

func (c *courseUseCase) FetchManyForEachCourseInAdmin(ctx context.Context, page string) ([]course_domain.CourseResponse, course_domain.DetailForManyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, detail, err := c.courseRepository.FetchManyForEachCourseInAdmin(ctx, page)
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}

	return course, detail, nil
}

func (c *courseUseCase) FindCourseIDByCourseNameInAdmin(ctx context.Context, courseName string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	course, err := c.courseRepository.FindCourseIDByCourseNameInAdmin(ctx, courseName)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return course, err
}

func (c *courseUseCase) CreateOneInAdmin(ctx context.Context, course *course_domain.Course) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	err := c.courseRepository.CreateOneInAdmin(ctx, course)

	if err != nil {
		return err
	}

	return nil
}

func (c *courseUseCase) UpdateOneInAdmin(ctx context.Context, course *course_domain.Course) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	data, err := c.courseRepository.UpdateOneInAdmin(ctx, course)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *courseUseCase) DeleteOneInAdmin(ctx context.Context, courseID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	err := c.courseRepository.DeleteOneInAdmin(ctx, courseID)
	if err != nil {
		return err
	}

	return err
}
