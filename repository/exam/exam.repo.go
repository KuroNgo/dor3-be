package exam_repository

import (
	exam_domain "clean-architecture/domain/exam"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type examRepository struct {
	database             mongo.Database
	collectionLesson     string
	collectionUnit       string
	collectionVocabulary string
	collectionExam       string
}

func NewExamRepository(db mongo.Database, collectionExam string, collectionLesson string, collectionUnit string, collectionVocabulary string) exam_domain.IExamRepository {
	return &examRepository{
		database:             db,
		collectionExam:       collectionExam,
		collectionLesson:     collectionLesson,
		collectionUnit:       collectionUnit,
		collectionVocabulary: collectionVocabulary,
	}
}

func (e *examRepository) FetchMany(ctx context.Context) (exam_domain.Response, error) {
	collectionExam := e.database.Collection(e.collectionExam)

	cursor, err := collectionExam.Find(ctx, bson.D{})
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

	err = cursor.All(ctx, &exams)
	examRes := exam_domain.Response{
		Exam: exams,
	}

	return examRes, err
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

func (e *examRepository) UpdateOne(ctx context.Context, examID string, exam exam_domain.Exam) error {
	collection := e.database.Collection(e.collectionExam)
	doc, err := internal.ToDoc(exam)
	objID, err := primitive.ObjectIDFromHex(examID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (e *examRepository) CreateOne(ctx context.Context, exam *exam_domain.Exam) error {
	collectionExam := e.database.Collection(e.collectionExam)
	collectionLesson := e.database.Collection(e.collectionLesson)
	collectionUnit := e.database.Collection(e.collectionUnit)
	collectionVocabulary := e.database.Collection(e.collectionVocabulary)

	filter := bson.M{"question": exam.Question}
	// check exists with CountDocuments
	count, err := collectionExam.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the lesson name did exist")
	}

	filterLessonID := bson.M{"_id": exam.LessonID}
	countLessonID, err := collectionLesson.CountDocuments(ctx, filterLessonID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the lesson ID do not exist")
	}

	filterUnitID := bson.M{"_id": exam.UnitID}
	countUnitID, err := collectionUnit.CountDocuments(ctx, filterUnitID)
	if err != nil {
		return err
	}
	if countUnitID == 0 {
		return errors.New("the unit ID do not exist")
	}

	filterVocabularyID := bson.M{"_id": exam.UnitID}
	countVocabularyID, err := collectionVocabulary.CountDocuments(ctx, filterVocabularyID)
	if err != nil {
		return err
	}
	if countVocabularyID == 0 {
		return errors.New("the vocabulary ID do not exist")
	}

	_, err = collectionLesson.InsertOne(ctx, exam)
	return nil
}

func (e *examRepository) UpdateCompleted(ctx context.Context, examID string, isComplete int) error {
	//TODO implement me
	panic("implement me")
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
