package course_repository

import (
	course_domain "clean-architecture/domain/course"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type courseRepository struct {
	database         *mongo.Database
	collectionCourse string
	collectionLesson string
}

func NewCourseRepository(db *mongo.Database, collectionCourse string, collectionLesson string) course_domain.ICourseRepository {
	return &courseRepository{
		database:         db,
		collectionCourse: collectionCourse,
		collectionLesson: collectionLesson,
	}
}

func (c *courseRepository) FetchMany(ctx context.Context) (course_domain.Response, error) {
	collectionCourse := c.database.Collection(c.collectionCourse)

	cursor, err := collectionCourse.Find(ctx, bson.D{})
	if err != nil {
		return course_domain.Response{}, err
	}

	var courses []course_domain.Course
	for cursor.Next(ctx) {
		var course course_domain.Course
		if err = cursor.Decode(&course); err != nil {
			return course_domain.Response{}, err
		}

		// Thêm lesson vào slice lessons
		courses = append(courses, course)
	}
	err = cursor.All(ctx, &courses)
	courseRes := course_domain.Response{
		Course: courses,
	}

	return courseRes, err
}

func (c *courseRepository) UpdateOne(ctx context.Context, course course_domain.Course) (*mongo.UpdateResult, error) {
	collectionCourse := c.database.Collection(c.collectionCourse)

	filter := bson.D{{Key: "_id", Value: course.Id}}
	update := bson.M{
		"$set": bson.M{
			"name":        course.Name,
			"description": course.Description,
		},
	}

	data, err := collectionCourse.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *courseRepository) CreateOne(ctx context.Context, course *course_domain.Course) error {
	collectionCourse := c.database.Collection(c.collectionCourse)

	filter := bson.M{"name": course.Name}
	// check exists with CountDocuments
	count, err := collectionCourse.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the course name did exist")
	}

	_, err = collectionCourse.InsertOne(ctx, course)
	return err
}

func (c *courseRepository) DeleteOne(ctx context.Context, courseID string) error {
	collectionCourse := c.database.Collection(c.collectionCourse)

	// Convert courseID string to ObjectID
	objID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return err
	}

	// Check if any lesson is associated with the course
	countFK, err := c.countLessonsByCourseID(ctx, courseID)
	if err != nil {
		return err
	}
	if countFK > 0 {
		return errors.New("the course cannot be deleted because it is associated with lessons")
	}

	// Delete the course
	filter := bson.M{"_id": objID}
	result, err := collectionCourse.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result == nil {
		return errors.New("the course was not found or already deleted")
	}

	return nil
}

// countLessonsByCourseID counts the number of lessons associated with a course.
func (c *courseRepository) countLessonsByCourseID(ctx context.Context, courseID string) (int64, error) {
	collectionLesson := c.database.Collection(c.collectionLesson)

	filter := bson.M{"course_id": courseID}
	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
