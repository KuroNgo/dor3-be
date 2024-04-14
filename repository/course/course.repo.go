package course_repository

import (
	course_domain "clean-architecture/domain/course"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (c *courseRepository) UpdateOne(ctx context.Context, courseID string, course course_domain.Course) error {
	collectionCourse := c.database.Collection(c.collectionCourse)
	doc, err := internal.ToDoc(course)
	objID, err := primitive.ObjectIDFromHex(courseID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collectionCourse.UpdateOne(ctx, filter, update)
	return err
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

func (c *courseRepository) UpsertOne(ctx context.Context, id string, course *course_domain.Course) (*course_domain.Response, error) {
	collectionCourse := c.database.Collection(c.collectionCourse)
	doc, err := internal.ToDoc(course)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionCourse.FindOneAndUpdate(ctx, query, update, opts)

	var updatedPost *course_domain.Response

	if err := res.Decode(&updatedPost); err != nil {
		return nil, errors.New("no post with that Id exists")
	}

	return updatedPost, nil
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
