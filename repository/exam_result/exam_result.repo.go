package exam_result_repository

import (
	exam_result_domain "clean-architecture/domain/exam_result"
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

type examResultRepository struct {
	database             *mongo.Database
	collectionExamResult string
	collectionAnswer     string
	collectionExam       string
}

func (e *examResultRepository) GetResultByID(ctx context.Context, id string) (exam_result_domain.ExamResult, error) {
	collectionResult := e.database.Collection(e.collectionExamResult)

	iD, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return exam_result_domain.ExamResult{}, err
	}

	filter := bson.M{"_id": iD}
	var res exam_result_domain.ExamResult
	err = collectionResult.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return exam_result_domain.ExamResult{}, err
	}
	return res, nil
}

func NewExamResultRepository(db *mongo.Database, collectionExamResult string, collectionExam string) exam_result_domain.IExamResultRepository {
	return &examResultRepository{
		database:             db,
		collectionExamResult: collectionExamResult,
		collectionExam:       collectionExam,
	}
}

func (e *examResultRepository) FetchManyByUserID(ctx context.Context, userID string) (exam_result_domain.Response, error) {
	collectionResult := e.database.Collection(e.collectionExamResult)

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return exam_result_domain.Response{}, err
	}

	filter := bson.M{"user_id": idUser}
	cursor, err := collectionResult.Find(ctx, filter)
	if err != nil {
		return exam_result_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var results []exam_result_domain.ExamResult
	score, _ := e.GetScoreByUser(ctx, userID)

	for cursor.Next(ctx) {
		var result exam_result_domain.ExamResult
		if err = cursor.Decode(&result); err != nil {
			return exam_result_domain.Response{}, err
		}

		result.UserID = idUser
		result.Score = score
		results = append(results, result)
	}
	averageScore, _ := e.GetAverageScoreByUser(ctx, userID)
	percentScore, _ := e.GetOverallPerformance(ctx, userID)

	resultRes := exam_result_domain.Response{
		ExamResult:   results,
		AverageScore: averageScore,
		Percentage:   percentScore,
	}

	return resultRes, nil
}

func (e *examResultRepository) GetResultsByUserIDAndExamID(ctx context.Context, userID string, examID string) (exam_result_domain.ExamResult, error) {
	collection := e.database.Collection(e.collectionExamResult)

	var result exam_result_domain.ExamResult

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return exam_result_domain.ExamResult{}, err
	}

	idExam, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return exam_result_domain.ExamResult{}, err
	}

	filter := bson.M{"user_id": idUser, "exam_id": idExam}

	err = collection.FindOne(ctx, filter).Decode(&result)
	return result, err
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

	count, err := collectionResult.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exam_result_domain.Response{}, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

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
		Page:       cal,
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
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var results []exam_result_domain.ExamResult

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var result exam_result_domain.ExamResult
			if err = cursor.Decode(&result); err != nil {
				return
			}

			result.ExamID = idExam
			results = append(results, result)
		}
	}()

	wg.Wait()

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

	_, err = collectionResult.InsertOne(ctx, examResult)
	return nil
}

func (e *examResultRepository) UpdateStatus(ctx context.Context, examID string, status int) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionExamResult)

	filter := bson.D{{Key: "exam_id", Value: examID}}
	update := bson.M{
		"$set": bson.M{
			"is_complete": status,
			"started_at":  time.Now(),
		},
	}

	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	return data, nil
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

func (e *examResultRepository) CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int {
	if totalQuestions == 0 {
		return 0
	}

	score := (correctAnswers * 10) / totalQuestions
	return score
}

func (e *examResultRepository) CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64 {
	if totalQuestions == 0 {
		return 0
	}

	percentage := float64(correctAnswers) / float64(totalQuestions) * 100
	return percentage
}

func (e *examResultRepository) GetScoreByUser(ctx context.Context, userID string) (int16, error) {
	// Lấy tất cả các kết quả của bài kiểm tra từ cơ sở dữ liệu MongoDB
	results, err := e.FetchManyByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Tính tổng điểm
	var totalScore int16 = 0
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, result := range results.ExamResult {
			totalScore += result.Score
		}
	}()

	wg.Wait()

	return totalScore, nil
}

func (e *examResultRepository) GetAverageScoreByUser(ctx context.Context, userID string) (float64, error) {
	results, err := e.FetchManyByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	var totalScore int16 = 0
	for _, result := range results.ExamResult {
		totalScore += result.Score
	}

	if len(results.ExamResult) > 0 {
		averageScore := float64(totalScore) / float64(len(results.ExamResult))
		return averageScore, nil
	}

	return 0, nil
}

func (e *examResultRepository) GetOverallPerformance(ctx context.Context, userID string) (float64, error) {
	averageScore, err := e.GetAverageScoreByUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	if averageScore == 0 {
		return 0, nil
	}

	overallPerformance := averageScore * 100

	return overallPerformance, nil
}
