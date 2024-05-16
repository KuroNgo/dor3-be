package exam_repository

import (
	exam_domain "clean-architecture/domain/exam"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type examRepository struct {
	database               *mongo.Database
	collectionLesson       string
	collectionUnit         string
	collectionExam         string
	collectionExamQuestion string

	examResponseCache map[string]exam_domain.DetailResponse
	examManyCache     map[string][]exam_domain.ExamResponse
	examOneCache      map[string]exam_domain.ExamResponse
	examCacheExpires  map[string]time.Time
	cacheMutex        sync.RWMutex
}

func NewExamRepository(db *mongo.Database, collectionExam string, collectionLesson string, collectionUnit string, collectionExamQuestion string) exam_domain.IExamRepository {
	return &examRepository{
		database:               db,
		collectionExam:         collectionExam,
		collectionLesson:       collectionLesson,
		collectionUnit:         collectionUnit,
		collectionExamQuestion: collectionExamQuestion,

		examResponseCache: make(map[string]exam_domain.DetailResponse),
		examManyCache:     make(map[string][]exam_domain.ExamResponse),
		examOneCache:      make(map[string]exam_domain.ExamResponse),
		examCacheExpires:  make(map[string]time.Time),
	}
}

func (e *examRepository) FetchMany(ctx context.Context, page string) ([]exam_domain.ExamResponse, exam_domain.DetailResponse, error) {
	//e.cacheMutex.RLock()
	//cacheData, found := e.examManyCache["exam"]
	//cachedResponseData, found := e.examResponseCache[page]
	//e.cacheMutex.RUnlock()
	//
	//if found {
	//	return cacheData, cachedResponseData, nil
	//}

	collectionExam := e.database.Collection(e.collectionExam)
	collectionUnit := e.database.Collection(e.collectionUnit)
	collectionLesson := e.database.Collection(e.collectionLesson)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))
	count, err := collectionExam.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	calCh := make(chan int64)
	go func() {
		defer close(calCh)
		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1 + 1
		}
	}()

	cursor, err := collectionExam.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var exams []exam_domain.ExamResponse

	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			var exam exam_domain.Exam
			if err = cursor.Decode(&exam); err != nil {
				return
			}
			countQuest := e.CountQuestion(ctx, exam.ID.Hex())

			var unit unit_domain.Unit
			filterUnit := bson.M{"_id": exam.UnitID}
			err = collectionUnit.FindOne(ctx, filterUnit).Decode(&unit)
			if err != nil {
				return
			}

			var lesson lesson_domain.Lesson
			filterLesson := bson.M{"_id": unit.LessonID}
			err = collectionLesson.FindOne(ctx, filterLesson).Decode(&lesson)
			if err != nil {
				return
			}

			var examRes exam_domain.ExamResponse
			examRes.ID = exam.ID
			examRes.Title = exam.Title
			examRes.Description = exam.Description
			examRes.Duration = exam.Duration
			examRes.CreatedAt = exam.CreatedAt
			examRes.UpdatedAt = exam.UpdatedAt
			examRes.WhoUpdates = exam.WhoUpdates
			examRes.Learner = exam.Learner
			examRes.IsComplete = exam.IsComplete
			examRes.CountQuestion = countQuest
			examRes.Unit = unit
			examRes.Lesson = lesson

			exams = append(exams, examRes)
		}
	}()

	internal.Wg.Wait()
	statisticsCh := make(chan exam_domain.Statistics)
	go func() {
		statistics, _ := e.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	cal := <-calCh
	detail := exam_domain.DetailResponse{
		Page:        cal,
		CurrentPage: pageNumber,
		CountExam:   count,
		Statistics:  statistics,
	}

	return exams, detail, nil
}

func (e *examRepository) FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]exam_domain.ExamResponse, exam_domain.DetailResponse, error) {
	//e.cacheMutex.RLock()
	//cacheData, found := e.examManyCache[unitID]
	//cachedResponseData, found := e.examResponseCache[page]
	//e.cacheMutex.RUnlock()
	//
	//if found {
	//	return cacheData, cachedResponseData, nil
	//}

	collectionExam := e.database.Collection(e.collectionExam)
	collectionUnit := e.database.Collection(e.collectionUnit)
	collectionLesson := e.database.Collection(e.collectionLesson)

	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		return nil, exam_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.D{{"_id", -1}})

	// Convert unitID to ObjectID
	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	// Count documents for pagination
	count, err := collectionExam.CountDocuments(ctx, bson.M{"unit_id": idUnit})
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	totalPages := (count + int64(perPage) - 1) / int64(perPage) // Calculate total pages

	// Query for exercises
	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionExam.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var exams []exam_domain.ExamResponse

	// Process each exercise
	for cursor.Next(ctx) {
		var exam exam_domain.Exam
		if err = cursor.Decode(&exam); err != nil {
			return nil, exam_domain.DetailResponse{}, err
		}

		// Fetch related data
		countQuest := e.CountQuestion(ctx, exam.ID.Hex())

		var unit unit_domain.Unit
		if err = collectionUnit.FindOne(ctx, bson.M{"_id": idUnit}).Decode(&unit); err != nil {
			return nil, exam_domain.DetailResponse{}, err
		}

		var lesson lesson_domain.Lesson
		if err = collectionLesson.FindOne(ctx, bson.M{"_id": unit.LessonID}).Decode(&lesson); err != nil {
			return nil, exam_domain.DetailResponse{}, err
		}

		var examRes exam_domain.ExamResponse
		examRes.ID = exam.ID
		examRes.Title = exam.Title
		examRes.Description = exam.Description
		examRes.Duration = exam.Duration
		examRes.CreatedAt = exam.CreatedAt
		examRes.UpdatedAt = exam.UpdatedAt
		examRes.WhoUpdates = exam.WhoUpdates
		examRes.Learner = exam.Learner
		examRes.IsComplete = exam.IsComplete
		examRes.CountQuestion = countQuest
		examRes.Unit = unit
		examRes.Lesson = lesson

		exams = append(exams, examRes)
	}

	if err = cursor.Err(); err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	detail := exam_domain.DetailResponse{
		CountExam:   count,
		Page:        totalPages,
		CurrentPage: pageNumber,
	}

	return exams, detail, nil
	//e.cacheMutex.Lock()
	//e.examManyCache[unitID] = exams
	//e.examResponseCache[page] = response
	//e.examCacheExpires[page] = time.Now().Add(5 * time.Minute)
	//e.examCacheExpires[unitID] = time.Now().Add(5 * time.Minute)
	//e.cacheMutex.Unlock()
}

func (e *examRepository) FetchOneByUnitID(ctx context.Context, unitID string) (exam_domain.ExamResponse, error) {
	//e.cacheMutex.RLock()
	//cacheData, found := e.examOneCache[unitID]
	//e.cacheMutex.RUnlock()

	//if found {
	//	return cacheData, nil
	//}

	collectionExam := e.database.Collection(e.collectionExam)
	collectionUnit := e.database.Collection(e.collectionUnit)
	collectionLesson := e.database.Collection(e.collectionLesson)

	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return exam_domain.ExamResponse{}, err
	}

	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionExam.Find(ctx, filter)
	if err != nil {
		return exam_domain.ExamResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var exams []exam_domain.ExamResponse
	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			var exam exam_domain.Exam
			if err = cursor.Decode(&exam); err != nil {
				return
			}

			// Fetch related data
			countQuest := e.CountQuestion(ctx, exam.ID.Hex())
			if err != nil {
				return
			}

			var unit unit_domain.Unit
			if err = collectionUnit.FindOne(ctx, bson.M{"_id": idUnit}).Decode(&unit); err != nil {
				return
			}

			var lesson lesson_domain.Lesson
			if err = collectionLesson.FindOne(ctx, bson.M{"_id": unit.LessonID}).Decode(&lesson); err != nil {
				return
			}

			var examRes exam_domain.ExamResponse
			examRes.ID = exam.ID
			examRes.Title = exam.Title
			examRes.Description = exam.Description
			examRes.Duration = exam.Duration
			examRes.CreatedAt = exam.CreatedAt
			examRes.UpdatedAt = exam.UpdatedAt
			examRes.WhoUpdates = exam.WhoUpdates
			examRes.Learner = exam.Learner
			examRes.IsComplete = exam.IsComplete
			examRes.CountQuestion = countQuest
			examRes.Unit = unit
			examRes.Lesson = lesson

			exams = append(exams, examRes)
		}
	}()
	internal.Wg.Wait()

	// Kiểm tra nếu danh sách exercises không rỗng
	if len(exams) == 0 {
		return exam_domain.ExamResponse{}, errors.New("no exercises found")
	}

	// Chọn một giá trị ngẫu nhiên từ danh sách exercises
	randomIndex := rand.Intn(len(exams))
	randomExam := exams[randomIndex]

	return randomExam, nil

	//e.cacheMutex.Lock()
	//e.examOneCache[unitID] = exam
	//e.examCacheExpires[unitID] = time.Now().Add(5 * time.Minute)
	//e.cacheMutex.Unlock()
}

func (e *examRepository) UpdateOne(ctx context.Context, exam *exam_domain.Exam) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionExam)

	filter := bson.M{"_id": exam.ID}
	update := bson.M{
		"$set": bson.M{
			"description": exam.Description,
			"title":       exam.Title,
			"duration":    exam.Duration,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (e *examRepository) CreateOne(ctx context.Context, exam *exam_domain.Exam) error {
	collectionExam := e.database.Collection(e.collectionExam)
	collectionLesson := e.database.Collection(e.collectionLesson)
	collectionUnit := e.database.Collection(e.collectionUnit)

	filterLessonID := bson.M{"_id": exam.LessonID}
	countLessonID, err := collectionLesson.CountDocuments(ctx, filterLessonID)
	if err != nil {
		return err
	}
	if countLessonID == 0 {
		return errors.New("the lesson ID does not exist")
	}

	filterUnitID := bson.M{"_id": exam.UnitID}
	countUnitID, err := collectionUnit.CountDocuments(ctx, filterUnitID)
	if err != nil {
		return err
	}
	if countUnitID == 0 {
		return errors.New("the unit ID does not exist")
	}

	_, err = collectionExam.InsertOne(ctx, exam)
	if err != nil {
		return err
	}

	return nil
}

func (e *examRepository) UpdateCompleted(ctx context.Context, exam *exam_domain.Exam) error {
	collection := e.database.Collection(e.collectionExam)

	filter := bson.D{{Key: "_id", Value: exam.ID}}
	update := bson.M{
		"$set": bson.M{
			"is_complete": exam.IsComplete,
			"update_at":   time.Now(),
			"learner":     exam.Learner,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (e *examRepository) DeleteOne(ctx context.Context, examID string) error {
	collectionExam := e.database.Collection(e.collectionExam)
	objID, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionExam.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam is removed`)
	}

	_, err = collectionExam.DeleteOne(ctx, filter)
	return err
}

func (e *examRepository) CountQuestion(ctx context.Context, examID string) int64 {
	collectionExamQuestion := e.database.Collection(e.collectionExamQuestion)

	idExam, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return 0
	}

	filter := bson.M{"exam_id": idExam}
	count, err := collectionExamQuestion.CountDocuments(ctx, filter)
	if err != nil {
		return 0
	}

	return count
}

func (e *examRepository) Statistics(ctx context.Context) (exam_domain.Statistics, error) {
	collectionExam := e.database.Collection(e.collectionExam)

	count, err := collectionExam.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exam_domain.Statistics{}, err
	}

	statistics := exam_domain.Statistics{
		Total: count,
	}
	return statistics, nil
}
