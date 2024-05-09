package course_repository

import (
	course_domain "clean-architecture/domain/course"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
	"time"
)

type courseRepository struct {
	database             *mongo.Database
	collectionCourse     string
	collectionLesson     string
	collectionUnit       string
	collectionVocabulary string

	courseResponseCache map[string]course_domain.DetailForManyResponse
	courseManyCache     map[string][]course_domain.CourseResponse
	courseOneCache      map[string]course_domain.CourseResponse
	courseCacheExpires  map[string]time.Time
	cacheMutex          sync.RWMutex
}

func NewCourseRepository(db *mongo.Database, collectionCourse string, collectionLesson string, collectionUnit string, collectionVocabulary string) course_domain.ICourseRepository {
	return &courseRepository{
		database:             db,
		collectionCourse:     collectionCourse,
		collectionLesson:     collectionLesson,
		collectionUnit:       collectionUnit,
		collectionVocabulary: collectionVocabulary,

		courseResponseCache: make(map[string]course_domain.DetailForManyResponse),
		courseManyCache:     make(map[string][]course_domain.CourseResponse),
		courseOneCache:      make(map[string]course_domain.CourseResponse),
		courseCacheExpires:  make(map[string]time.Time),
	}
}

func (c *courseRepository) FetchByID(ctx context.Context, courseID string) (course_domain.CourseResponse, error) {
	c.cacheMutex.RLock()
	cachedData, found := c.courseOneCache[courseID]
	c.cacheMutex.RUnlock()

	if found {
		return cachedData, nil
	}

	collectionCourse := c.database.Collection(c.collectionCourse)
	idCourse, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return course_domain.CourseResponse{}, err
	}
	filter := bson.M{"_id": idCourse}

	var course course_domain.CourseResponse
	err = collectionCourse.FindOne(ctx, filter).Decode(&course)
	if err != nil {
		return course_domain.CourseResponse{}, err
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

	course.CountVocabulary = countVocab
	course.CountLesson = countLesson

	c.cacheMutex.Lock()
	c.courseOneCache[courseID] = course
	c.courseCacheExpires[courseID] = time.Now().Add(5 * time.Minute)
	c.cacheMutex.Unlock()

	return course, nil
}

func (c *courseRepository) FetchManyForEachCourse(ctx context.Context, page string) ([]course_domain.CourseResponse, course_domain.DetailForManyResponse, error) {
	c.cacheMutex.RLock()
	cachedData, found := c.courseManyCache["course"]
	cachedResponseData, found := c.courseResponseCache[page]
	c.cacheMutex.RUnlock()

	if found {
		return cachedData, cachedResponseData, nil
	}

	collectionCourse := c.database.Collection(c.collectionCourse)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, errors.New("invalid page number")
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	count, err := collectionCourse.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}

	totalPages := (count + int64(perPage) - 1) / int64(perPage)
	cursor, err := collectionCourse.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}()

	var courses []course_domain.CourseResponse
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var course course_domain.CourseResponse
			if err := cursor.Decode(&course); err != nil {
				return
			}

			wg.Add(1)
			go func(course course_domain.CourseResponse) {
				defer wg.Done()

				countLesson, err := c.countLessonsByCourseID(ctx, course.Id)
				if err != nil {
					return
				}

				countVocab, err := c.countVocabularyByCourseID(ctx, course.Id)
				if err != nil {
					return
				}

				course.CountVocabulary = countVocab
				course.CountLesson = countLesson

				courses = append(courses, course)
			}(course)
		}
	}()

	wg.Wait()

	detail := course_domain.DetailForManyResponse{
		CountCourse: count,
		Page:        totalPages,
		CurrentPage: pageNumber,
	}

	c.cacheMutex.Lock()
	c.courseManyCache["course"] = courses
	c.courseResponseCache[page] = detail
	c.courseCacheExpires["course"] = time.Now().Add(5 * time.Minute)
	c.courseCacheExpires[page] = time.Now().Add(5 * time.Minute)
	c.cacheMutex.Unlock()

	return courses, detail, nil
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

	var mu sync.Mutex // Mutex để bảo vệ courses

	mu.Lock()
	data, err := collectionCourse.UpdateOne(ctx, filter, &update)
	mu.Unlock()

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
