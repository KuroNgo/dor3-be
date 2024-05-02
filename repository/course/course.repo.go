package course_repository

import (
	course_domain "clean-architecture/domain/course"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

var (
	wg sync.WaitGroup
)

type courseRepository struct {
	database             *mongo.Database
	collectionCourse     string
	collectionLesson     string
	collectionUnit       string
	collectionVocabulary string
}

func NewCourseRepository(db *mongo.Database, collectionCourse string, collectionLesson string, collectionUnit string, collectionVocabulary string) course_domain.ICourseRepository {
	return &courseRepository{
		database:             db,
		collectionCourse:     collectionCourse,
		collectionLesson:     collectionLesson,
		collectionUnit:       collectionUnit,
		collectionVocabulary: collectionVocabulary,
	}
}

func (c *courseRepository) FetchManyForEachCourse(ctx context.Context) ([]course_domain.CourseResponse, error) {
	collectionCourse := c.database.Collection(c.collectionCourse)
	cursor, err := collectionCourse.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var courses []course_domain.CourseResponse
	for cursor.Next(ctx) {
		var course course_domain.Course
		if err := cursor.Decode(&course); err != nil {
			return nil, err
		}

		countLessonCh := make(chan int32)
		countVocabularyCh := make(chan int32)

		go func() {
			defer close(countLessonCh)
			countLesson, err := c.countLessonsByCourseID(ctx, course.Id)
			if err != nil {
				return
			}
			countLessonCh <- countLesson
		}()

		go func() {
			defer close(countVocabularyCh)
			countVocabulary, err := c.countVocabularyByCourseID(ctx, course.Id)
			if err != nil {
				return
			}
			countVocabularyCh <- countVocabulary
		}()

		countLesson := <-countLessonCh
		countVocab := <-countVocabularyCh

		courseRes := course_domain.CourseResponse{
			Id:              course.Id,
			Name:            course.Name,
			Description:     course.Description,
			CreatedAt:       course.CreatedAt,
			UpdatedAt:       course.UpdatedAt,
			WhoUpdated:      course.WhoUpdated,
			CountLesson:     countLesson,
			CountVocabulary: countVocab,
		}

		courses = append(courses, courseRes)
	}

	return courses, nil
}

func (c *courseRepository) UpdateOne(ctx context.Context, course *course_domain.Course) (*mongo.UpdateResult, error) {
	collectionCourse := c.database.Collection(c.collectionCourse)

	filter := bson.D{{Key: "_id", Value: course.Id}}
	update := bson.M{
		"$set": bson.M{
			"name":        course.Name,
			"description": course.Description,
			"updated_at":  course.UpdatedAt,
			"who_updated": course.WhoUpdated,
		},
	}

	data, err := collectionCourse.UpdateOne(ctx, filter, &update)
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

func (c *courseRepository) StatisticCourse(ctx context.Context) int64 {
	collectionCourse := c.database.Collection(c.collectionCourse)

	filter := bson.D{}

	count, err := collectionCourse.CountDocuments(ctx, filter)
	if err != nil {
		return 0
	}

	return count
}

func (c *courseRepository) DeleteOne(ctx context.Context, courseID string) error {
	collectionCourse := c.database.Collection(c.collectionCourse)

	// Default the Course for iT cannot delete
	objID2, err := primitive.ObjectIDFromHex("660b8a0c2aef1f3a28265523")
	if err != nil {
		return err
	}
	countIn, err := collectionCourse.CountDocuments(ctx, objID2)
	if countIn > 0 {
		return errors.New("the course cannot be deleted")
	}

	// Convert courseID string to ObjectID
	objID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return err
	}

	// Check if any lesson is associated with the course
	countFK, err := c.countLessonsByCourseID(ctx, objID)
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
func (c *courseRepository) countLessonsByCourseID(ctx context.Context, courseID primitive.ObjectID) (int32, error) {
	collectionLesson := c.database.Collection(c.collectionLesson)

	filter := bson.M{"course_id": courseID}
	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

func (c *courseRepository) countVocabularyByCourseID(ctx context.Context, courseID primitive.ObjectID) (int32, error) {
	collectionVocabulary := c.database.Collection(c.collectionVocabulary)

	// Phần pipeline aggregation để nối các collection và đếm số lượng từ vựng
	pipeline := []bson.M{
		// Nối collection Vocabulary với collection Unit
		{"$lookup": bson.M{
			"from":         "unit",
			"localField":   "unit_id",
			"foreignField": "_id",
			"as":           "unit",
		}},
		{"$unwind": "$unit"},
		// Nối collection Unit với collection Lesson
		{"$lookup": bson.M{
			"from":         "lesson",
			"localField":   "unit.lesson_id",
			"foreignField": "_id",
			"as":           "lesson",
		}},
		{"$unwind": "$lesson"},
		// Lọc các từ vựng thuộc về khóa học được cung cấp
		{"$match": bson.M{"lesson.course_id": courseID}},
		// Đếm số lượng từ vựng
		{"$count": "totalVocabulary"},
	}

	// Thực hiện aggregation
	cursor, err := collectionVocabulary.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var result struct {
		TotalVocabulary int32 `bson:"totalVocabulary"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.TotalVocabulary, nil
}
