package exam_repository

import (
	exam_domain "clean-architecture/domain/exam"
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
	collectionVocabulary   string
}

func NewExamRepository(db *mongo.Database, collectionExam string, collectionLesson string, collectionUnit string, collectionExamQuestion string, collectionVocabulary string) exam_domain.IExamRepository {
	return &examRepository{
		database:               db,
		collectionExam:         collectionExam,
		collectionLesson:       collectionLesson,
		collectionUnit:         collectionUnit,
		collectionExamQuestion: collectionExamQuestion,
		collectionVocabulary:   collectionVocabulary,
	}
}

var (
	wg sync.WaitGroup
)

func (e *examRepository) FetchMany(ctx context.Context, page string) ([]exam_domain.Exam, exam_domain.DetailResponse, error) {
	collectionExam := e.database.Collection(e.collectionExam)

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

	var exams []exam_domain.Exam
	for cursor.Next(ctx) {
		var exam exam_domain.Exam
		if err = cursor.Decode(&exam); err != nil {
			return nil, exam_domain.DetailResponse{}, err
		}

		exams = append(exams, exam)
	}
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
		CountExam:   int64(len(exams)),
		Statistics:  statistics,
	}

	return exams, detail, nil
}

func (e *examRepository) FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]exam_domain.Exam, exam_domain.DetailResponse, error) {
	collectionExam := e.database.Collection(e.collectionExam)

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

	var exams []exam_domain.Exam

	// Process each exercise
	for cursor.Next(ctx) {
		var exam exam_domain.Exam
		if err = cursor.Decode(&exam); err != nil {
			return nil, exam_domain.DetailResponse{}, err
		}

		exams = append(exams, exam)
	}

	if err = cursor.Err(); err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	statisticsCh := make(chan exam_domain.Statistics)
	go func() {
		statistics, _ := e.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	detail := exam_domain.DetailResponse{
		CountExam:   int64(len(exams)),
		Page:        totalPages,
		Statistics:  statistics,
		CurrentPage: pageNumber,
	}

	return exams, detail, nil
}

func (e *examRepository) FetchExamByID(ctx context.Context, id string) (exam_domain.Exam, error) {
	collectionExam := e.database.Collection(e.collectionExam)

	idExam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return exam_domain.Exam{}, err
	}

	var exam exam_domain.Exam
	filter := bson.M{"_id": idExam}
	err = collectionExam.FindOne(ctx, filter).Decode(&exam)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return exam_domain.Exam{}, errors.New("exam not found")
		}
		return exam_domain.Exam{}, err
	}

	return exam, nil
}

func (e *examRepository) FetchOneByUnitID(ctx context.Context, unitID string) (exam_domain.Exam, error) {
	collectionExam := e.database.Collection(e.collectionExam)

	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return exam_domain.Exam{}, err
	}

	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionExam.Find(ctx, filter)
	if err != nil {
		return exam_domain.Exam{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var exams []exam_domain.Exam

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var exam exam_domain.Exam
			if err = cursor.Decode(&exam); err != nil {
				return
			}

			exams = append(exams, exam)

		}
	}()
	wg.Wait()

	// Kiểm tra nếu danh sách exercises không rỗng
	if len(exams) == 0 {
		return exam_domain.Exam{}, errors.New("no exercises found")
	}

	// Chọn một giá trị ngẫu nhiên từ danh sách exercises
	randomIndex := rand.Intn(len(exams))
	randomExam := exams[randomIndex]

	return randomExam, nil
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

	filterUnitID := bson.M{"_id": exam.UnitID, "lesson_id": exam.LessonID}
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

func (e *examRepository) UpdateCompleted(ctx context.Context, exam *exam_domain.Exam) error {
	collection := e.database.Collection(e.collectionExam)

	filter := bson.D{{Key: "_id", Value: exam.ID}}
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
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
		return errors.New(`exam is removed or have not exist`)
	}

	_, err = collectionExam.DeleteOne(ctx, filter)
	return err
}

func (e *examRepository) countQuestion(ctx context.Context, examID string) int64 {
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
