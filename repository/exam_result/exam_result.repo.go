package exam_result

import (
	exam_result_domain "clean-architecture/domain/exam_result"
	"clean-architecture/infrastructor/mongo"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type examResultRepository struct {
	database             mongo.Database
	collectionExamResult string
	collectionExam       string
}

func NewExamRepository(db mongo.Database, collectionExamResult string, collectionExam string) exam_result_domain.IExamOptionsRepository {
	return &examResultRepository{
		database:             db,
		collectionExamResult: collectionExamResult,
		collectionExam:       collectionExam,
	}
}

func (e *examResultRepository) FetchMany(ctx context.Context, page string) (exam_result_domain.Response, error) {
	collectionResult := e.database.Collection(e.collectionExamResult)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return exam_result_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	cursor, err := collectionResult.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return exam_result_domain.Response{}, err
	}

	var results []exam_result_domain.ExamResult
	for cursor.Next(ctx) {
		var result exam_result_domain.ExamResult
		if err = cursor.Decode(&result); err != nil {
			return exam_result_domain.Response{}, err
		}

		results = append(results, result)
	}

	resultRes := exam_result_domain.Response{
		ExamResult: results,
	}

	return resultRes, nil
}

func (e *examResultRepository) FetchManyByExamID(ctx context.Context, examID string) (exam_result_domain.Response, error) {
	collectionResult := e.database.Collection(e.collectionExamResult)

	idExam, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return exam_result_domain.Response{}, err
	}

	filter := bson.M{"exam_id": idExam}
	cursor, err := collectionResult.Find(ctx, filter)
	if err != nil {
		return exam_result_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var results []exam_result_domain.ExamResult
	for cursor.Next(ctx) {
		var result exam_result_domain.ExamResult
		if err = cursor.Decode(&result); err != nil {
			return exam_result_domain.Response{}, err
		}

		result.ExamID = idExam
		results = append(results, result)
	}

	resultRes := exam_result_domain.Response{
		ExamResult: results,
	}
	return resultRes, nil
}

func (e *examResultRepository) CreateOne(ctx context.Context, examResult *exam_result_domain.ExamResult) error {
	collectionResult := e.database.Collection(e.collectionExamResult)
	collectionExam := e.database.Collection(e.collectionExam)

	filterExamID := bson.M{"exam_id": examResult.ExamID}
	countLessonID, err := collectionExam.CountDocuments(ctx, filterExamID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the examID do not exist")
	}

	_, err = collectionResult.InsertOne(ctx, collectionResult)
	return nil
}

func (e *examResultRepository) DeleteOne(ctx context.Context, examResultID string) error {
	collectionResult := e.database.Collection(e.collectionExamResult)

	objID, err := primitive.ObjectIDFromHex(examResultID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionResult.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam answer is removed`)
	}

	_, err = collectionResult.DeleteOne(ctx, filter)
	return err
}
