package lesson_repository

import (
	lesson_domain "clean-architecture/domain/lesson"
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

type lessonRepository struct {
	database             *mongo.Database
	collectionLesson     string
	collectionCourse     string
	collectionUnit       string
	collectionVocabulary string

	lessonResponseCache map[string]lesson_domain.DetailResponse
	lessonManyCache     map[string][]lesson_domain.LessonResponse
	lessonOneCache      map[string]lesson_domain.LessonResponse
	lessonCacheExpires  map[string]time.Time
	cacheMutex          sync.RWMutex
}

func NewLessonRepository(db *mongo.Database, collectionLesson string, collectionCourse string, collectionUnit string, collectionVocabulary string) lesson_domain.ILessonRepository {
	return &lessonRepository{
		database:             db,
		collectionLesson:     collectionLesson,
		collectionCourse:     collectionCourse,
		collectionUnit:       collectionUnit,
		collectionVocabulary: collectionVocabulary,

		lessonResponseCache: make(map[string]lesson_domain.DetailResponse),
		lessonManyCache:     make(map[string][]lesson_domain.LessonResponse),
		lessonOneCache:      make(map[string]lesson_domain.LessonResponse),
		lessonCacheExpires:  make(map[string]time.Time),
	}
}

func (l *lessonRepository) FetchByID(ctx context.Context, lessonID string) (lesson_domain.LessonResponse, error) {
	//l.cacheMutex.RLock()
	//cachedData, found := l.lessonOneCache[lessonID]
	//l.cacheMutex.RUnlock()
	//if found {
	//	return cachedData, nil
	//}

	collectionLesson := l.database.Collection(l.collectionLesson)

	idLesson, err := primitive.ObjectIDFromHex(lessonID)
	if err != nil {
		return lesson_domain.LessonResponse{}, err
	}

	filter := bson.M{"_id": idLesson}

	var lesson lesson_domain.LessonResponse
	err = collectionLesson.FindOne(ctx, filter).Decode(&lesson)
	if err != nil {
		return lesson_domain.LessonResponse{}, err
	}

	countUnitCh := make(chan int32)
	countVocabularyCh := make(chan int32)

	go func() {
		defer close(countVocabularyCh)
		countVocabulary, err := l.countVocabularyByLessonID(ctx, lesson.ID)
		if err != nil {
			return
		}
		countVocabularyCh <- countVocabulary
	}()

	go func() {
		defer close(countUnitCh)
		countUnit, err := l.countUnitsByLessonsID(ctx, lesson.ID)
		if err != nil {
			return
		}
		countUnitCh <- countUnit
	}()

	countUnit := <-countUnitCh
	countVocabulary := <-countVocabularyCh
	lesson.CountVocabulary = countVocabulary
	lesson.CountUnit = countUnit

	//l.cacheMutex.Lock()
	//l.lessonOneCache[lessonID] = lesson
	//l.lessonCacheExpires[lessonID] = time.Now().Add(5 * time.Minute)
	//l.cacheMutex.Unlock()

	return lesson, nil
}

func (l *lessonRepository) FetchByIdCourse(ctx context.Context, idCourse string, page string) ([]lesson_domain.LessonResponse, lesson_domain.DetailResponse, error) {
	//l.cacheMutex.RLock()
	//cachedData, found := l.lessonManyCache[idCourse]
	//cachedResponseData, found := l.lessonResponseCache[page]
	//l.cacheMutex.RUnlock()
	//if found {
	//	return cachedData, cachedResponseData, nil
	//}

	collectionLesson := l.database.Collection(l.collectionLesson)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	calCh := make(chan int64)

	go func() {
		count, err := collectionLesson.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1 + 1
		}
	}()

	idCourse2, err := primitive.ObjectIDFromHex(idCourse)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}

	filter := bson.M{"course_id": idCourse2}

	cursor, err := collectionLesson.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var lessons []lesson_domain.LessonResponse

	for cursor.Next(ctx) {
		var lesson lesson_domain.LessonResponse
		if err = cursor.Decode(&lesson); err != nil {
			return nil, lesson_domain.DetailResponse{}, err
		}

		countUnitCh := make(chan int32)
		go func() {
			defer close(countUnitCh)
			// Lấy thông tin liên quan cho mỗi chủ đề
			countUnit, err := l.countUnitsByLessonsID(ctx, lesson.ID)
			if err != nil {
				return
			}

			countUnitCh <- countUnit
		}()

		countVocabularyCh := make(chan int32)
		go func() {
			defer close(countVocabularyCh)
			countVocabulary, err := l.countVocabularyByLessonID(ctx, lesson.ID)
			if err != nil {
				return
			}

			countVocabularyCh <- countVocabulary
		}()

		countUnit := <-countUnitCh

		countVocabulary := <-countVocabularyCh

		// Gắn CourseID vào bài học
		lesson.CourseID = idCourse2
		lesson.CountVocabulary = countVocabulary
		lesson.CountUnit = countUnit

		lessons = append(lessons, lesson)
	}

	cal := <-calCh

	response := lesson_domain.DetailResponse{
		Page:        cal,
		CurrentPage: pageNumber,
	}

	//l.cacheMutex.Lock()
	//l.lessonManyCache[idCourse] = lessons
	//l.lessonResponseCache[page] = response
	//l.lessonCacheExpires[idCourse] = time.Now().Add(5 * time.Minute)
	//l.lessonCacheExpires[page] = time.Now().Add(5 * time.Minute)
	//l.cacheMutex.Unlock()

	return lessons, response, nil
}

func (l *lessonRepository) FindCourseIDByCourseName(ctx context.Context, courseName string) (primitive.ObjectID, error) {
	collectionCourse := l.database.Collection(l.collectionCourse)

	filter := bson.M{"name": courseName}
	var data struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	err := collectionCourse.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return data.Id, nil
}

func (l *lessonRepository) UpdateComplete(ctx context.Context, lessonID string, lesson lesson_domain.Lesson) error {
	collection := l.database.Collection(l.collectionUnit)

	filter := bson.D{{Key: "_id", Value: lessonID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "is_complete", Value: lesson.IsCompleted},
		{Key: "who_updates", Value: lesson.WhoUpdates},
	}}}

	_, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return err
	}
	return nil
}

func (l *lessonRepository) FetchMany(ctx context.Context, page string) ([]lesson_domain.LessonResponse, lesson_domain.DetailResponse, error) {
	//l.cacheMutex.RLock()
	//cachedData, found := l.lessonManyCache[page]
	//cachedResponseData, found := l.lessonResponseCache[page]
	//l.cacheMutex.RUnlock()
	//if found {
	//	return cachedData, cachedResponseData, nil
	//}

	collectionLesson := l.database.Collection(l.collectionLesson)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 10
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	calCh := make(chan int64)

	go func() {
		defer close(calCh)
		count, err := collectionLesson.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1 + 1
		}
	}()

	cursor, err := collectionLesson.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var lessons []lesson_domain.LessonResponse
	for cursor.Next(ctx) {
		var lesson lesson_domain.LessonResponse
		if err = cursor.Decode(&lesson); err != nil {
			return nil, lesson_domain.DetailResponse{}, err
		}

		countUnitCh := make(chan int32)
		go func() {
			defer close(countUnitCh)
			// Lấy thông tin liên quan cho mỗi chủ đề
			countUnit, err := l.countUnitsByLessonsID(ctx, lesson.ID)
			if err != nil {
				return
			}

			countUnitCh <- countUnit
		}()

		countVocabularyCh := make(chan int32)
		go func() {
			defer close(countVocabularyCh)
			countVocabulary, err := l.countVocabularyByLessonID(ctx, lesson.ID)
			if err != nil {
				return
			}

			countVocabularyCh <- countVocabulary
		}()

		countUnit := <-countUnitCh
		countVocabulary := <-countVocabularyCh

		lesson.CountUnit = countUnit
		lesson.CountVocabulary = countVocabulary

		// Thêm lesson vào slice lessons
		lessons = append(lessons, lesson)
	}

	cal := <-calCh
	response := lesson_domain.DetailResponse{
		Page:        cal,
		CurrentPage: pageNumber,
	}

	//l.cacheMutex.Lock()
	//l.lessonManyCache[page] = lessons
	//l.lessonResponseCache[page] = response
	//l.lessonCacheExpires[page] = time.Now().Add(5 * time.Minute)
	//l.cacheMutex.Unlock()

	return lessons, response, err
}

func (l *lessonRepository) UpdateOne(ctx context.Context, lesson *lesson_domain.Lesson) (*mongo.UpdateResult, error) {
	collection := l.database.Collection(l.collectionLesson)

	filter := bson.M{"_id": lesson.ID}

	update := bson.M{
		"$set": bson.M{
			"name":    lesson.Name,
			"content": lesson.Content,
			"image":   lesson.ImageURL,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (l *lessonRepository) CreateOne(ctx context.Context, lesson *lesson_domain.Lesson) error {
	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionCourse := l.database.Collection(l.collectionCourse)

	filter := bson.M{"name": lesson.Name}
	// check exists with CountDocuments
	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the lesson name did exist")
	}

	filterReference := bson.M{"_id": lesson.CourseID}
	countParent, err := collectionCourse.CountDocuments(ctx, filterReference)
	if err != nil {
		return err
	}

	if countParent == 0 {
		return errors.New("the course ID do not exist")
	}

	_, err = collectionLesson.InsertOne(ctx, lesson)
	return nil
}

func (l *lessonRepository) CreateOneByNameCourse(ctx context.Context, lesson *lesson_domain.Lesson) error {
	collectionLesson := l.database.Collection(l.collectionLesson)

	filter := bson.M{"name": lesson.Name}
	// check exists with CountDocuments
	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the lesson name did exist")
	}

	_, err = collectionLesson.InsertOne(ctx, lesson)
	return nil
}

func (l *lessonRepository) DeleteOne(ctx context.Context, lessonID string) error {
	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionUnit := l.database.Collection(l.collectionUnit)

	objID, err := primitive.ObjectIDFromHex(lessonID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	filterInUnit := bson.M{"lesson_id": objID}

	exist, err := collectionUnit.CountDocuments(ctx, filterInUnit)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New(`lesson can not remove`)
	}

	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`lesson is removed`)
	}

	_, err = collectionLesson.DeleteOne(ctx, filter)
	return err
}

// countLessonsByCourseID counts the number of lessons associated with a course.
func (l *lessonRepository) countUnitsByLessonsID(ctx context.Context, lessonID primitive.ObjectID) (int32, error) {
	collectionUnit := l.database.Collection(l.collectionUnit)

	filter := bson.M{"lesson_id": lessonID}
	count, err := collectionUnit.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

func (l *lessonRepository) countVocabularyByLessonID(ctx context.Context, lessonID primitive.ObjectID) (int32, error) {
	collectionVocabulary := l.database.Collection(l.collectionVocabulary)

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

		// Lọc các từ vựng thuộc về khóa học được cung cấp
		{"$match": bson.M{"unit.lesson_id": lessonID}},
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
