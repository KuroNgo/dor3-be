package course_repository

import (
	course_domain "clean-architecture/domain/course"
	"clean-architecture/internal"
	"clean-architecture/internal/cache"
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

var (
	courseCache  = cache.NewTTL[string, course_domain.CourseResponse]()
	coursesCache = cache.NewTTL[string, []course_domain.CourseResponse]()
	detailCache  = cache.NewTTL[string, course_domain.DetailForManyResponse]()

	wg sync.WaitGroup
	mu sync.Mutex

	statisticsCh = make(chan course_domain.Statistics)
)

func (c *courseRepository) FetchByID(ctx context.Context, courseID string) (course_domain.CourseResponse, error) {
	// implement channel
	courseCh := make(chan course_domain.CourseResponse)
	wg.Add(1)

	// a goroutine do check data value in cache
	go func() {
		defer wg.Done()
		value, found := courseCache.Get(courseID)
		if found {
			courseCh <- value
			return
		}
	}()

	// to prevent panic "send to channel close", implement goroutine
	go func() {
		defer close(courseCh)
		wg.Wait()
	}()

	courseData := <-courseCh
	if !internal.IsZeroValue(courseData) {
		return courseData, nil
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

	courseCache.Set(courseID, course, 5*time.Minute)

	return course, nil
}

func (c *courseRepository) FetchManyForEachCourse(ctx context.Context, page string) ([]course_domain.CourseResponse, course_domain.DetailForManyResponse, error) {
	// buffer channel with increase performance
	// note: use buffer have target
	coursesCh := make(chan []course_domain.CourseResponse, 5)
	detailCh := make(chan course_domain.DetailForManyResponse, 1)
	wg.Add(2)
	go func() {
		defer wg.Done()
		data, found := coursesCache.Get(page)
		if found {
			coursesCh <- data
			return
		}
	}()
	go func() {
		defer wg.Done()
		detailData, foundDetail := detailCache.Get("detail")
		if foundDetail {
			detailCh <- detailData
			return
		}
	}()

	go func() {
		defer close(coursesCh)
		defer close(detailCh)
		wg.Wait()
	}()

	courseData := <-coursesCh
	responseData := <-detailCh
	if !internal.IsZeroValue(courseData) && !internal.IsZeroValue(responseData) {
		return courseData, responseData, nil
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var course course_domain.CourseResponse
			if err := cursor.Decode(&course); err != nil {
				return
			}

			wg.Add(1)
			go func(course2 course_domain.CourseResponse) {
				defer wg.Done()
				countLesson, err := c.countLessonsByCourseID(ctx, course2.Id)
				if err != nil {
					return
				}

				countVocab, err := c.countVocabularyByCourseID(ctx, course2.Id)
				if err != nil {
					return
				}
				course2.CountVocabulary = countVocab
				course2.CountLesson = countLesson
				courses = append(courses, course2)
			}(course)
		}
	}()

	wg.Wait()

	go func() {
		statistics, _ := c.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	detail := course_domain.DetailForManyResponse{
		CountCourse: count,
		Page:        totalPages,
		Statistics:  statistics,
		CurrentPage: pageNumber,
	}

	coursesCache.Set(page, courses, 5*time.Minute)
	detailCache.Set("detail", detail, 5*time.Minute)

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

	mu.Lock()
	data, err := collectionCourse.UpdateOne(ctx, filter, &update)
	mu.Unlock()
	if err != nil {
		return nil, err
	}

	// Clear data value in cache memory for courses
	wg.Add(2)
	go func() {
		defer wg.Done()
		coursesCache.Clear()
	}()

	// clear data value with id courseID in cache
	go func() {
		defer wg.Done()
		courseCache.Remove(course.Id.Hex())
	}()
	wg.Wait()

	return data, nil
}

func (c *courseRepository) CreateOne(ctx context.Context, course *course_domain.Course) error {
	collectionCourse := c.database.Collection(c.collectionCourse)

	filter := bson.M{"name": course.Name}
	count, err := collectionCourse.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the course name already exists")
	}

	mu.Lock()
	_, err = collectionCourse.InsertOne(ctx, course)
	mu.Unlock()
	if err != nil {
		return err
	}

	// Clear data value in cache memory
	wg.Add(2)
	go func() {
		defer wg.Done()
		coursesCache.Clear()
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		detailCache.Clear()
	}()

	wg.Wait()

	return nil
}

func (c *courseRepository) DeleteOne(ctx context.Context, courseID string) error {
	collectionCourse := c.database.Collection(c.collectionCourse)

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
	mu.Lock()
	_, err = collectionCourse.DeleteOne(ctx, filter)
	mu.Unlock()
	if err != nil {
		return err
	}

	// clear data value with courseID in cache
	wg.Add(2)
	go func() {
		defer wg.Done()
		courseCache.Remove(courseID)
	}()

	// clear data value with detail in cache due to decrease num
	go func() {
		defer wg.Done()
		detailCache.Clear()
	}()
	wg.Wait()

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

func (c *courseRepository) Statistics(ctx context.Context) (course_domain.Statistics, error) {
	collectionCourse := c.database.Collection(c.collectionCourse)

	count, err := collectionCourse.CountDocuments(ctx, bson.D{})
	if err != nil {
		return course_domain.Statistics{}, err
	}

	statistics := course_domain.Statistics{
		Total: count,
	}
	return statistics, nil
}
