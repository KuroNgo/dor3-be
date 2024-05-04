package exam_repository

import (
	exam_domain "clean-architecture/domain/exam"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type examRepository struct {
	database               *mongo.Database
	collectionLesson       string
	collectionUnit         string
	collectionExam         string
	collectionExamQuestion string
}

func (e *examRepository) FetchOneByUnitID(ctx context.Context, unitID string) (exam_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func NewExamRepository(db *mongo.Database, collectionExam string, collectionLesson string, collectionUnit string, collectionExamQuestion string) exam_domain.IExamRepository {
	return &examRepository{
		database:               db,
		collectionExam:         collectionExam,
		collectionLesson:       collectionLesson,
		collectionUnit:         collectionUnit,
		collectionExamQuestion: collectionExamQuestion,
	}
}

func (e *examRepository) FetchMany(ctx context.Context, page string) (exam_domain.Response, error) {
	collectionExam := e.database.Collection(e.collectionExam)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return exam_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))
	count, err := collectionExam.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exam_domain.Response{}, err
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
		return exam_domain.Response{}, err
	}

	var exams []exam_domain.Exam
	for cursor.Next(ctx) {
		var exam exam_domain.Exam
		if err = cursor.Decode(&exam); err != nil {
			return exam_domain.Response{}, err
		}

		// Thêm lesson vào slice lessons
		exams = append(exams, exam)
	}

	cal := <-calCh
	examRes := exam_domain.Response{
		Page:      cal,
		Exam:      exams,
		CountExam: count,
	}

	return examRes, nil
}

func (e *examRepository) FetchManyByUnitID(ctx context.Context, unitID string) (exam_domain.Response, error) {
	collectionExam := e.database.Collection(e.collectionExam)

	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return exam_domain.Response{}, err
	}

	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionExam.Find(ctx, filter)
	if err != nil {
		return exam_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var exams []exam_domain.Exam
	for cursor.Next(ctx) {
		var exam exam_domain.Exam
		if err = cursor.Decode(&exam); err != nil {
			return exam_domain.Response{}, err
		}

		// Gắn CourseID vào bài học
		exam.UnitID = idUnit
		exams = append(exams, exam)
	}

	response := exam_domain.Response{
		Exam: exams,
	}
	return response, nil
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
