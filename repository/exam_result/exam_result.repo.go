package exam_result

import (
	exam_result_domain "clean-architecture/domain/exam_result"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type examResultRepository struct {
	database             mongo.Database
	collectionExamResult string
	collectionExam       string
}

func NewExamResultRepository(db mongo.Database, collectionExamResult string, collectionExam string) exam_result_domain.IExamResultRepository {
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

	_, err = collectionResult.InsertOne(ctx, examResult)
	return nil
}

func (e *examResultRepository) UpdateStatus(ctx context.Context, examResultID string, status int) error {
	collection := e.database.Collection(e.collectionExamResult)
	objID, err := primitive.ObjectIDFromHex(examResultID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

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
	for _, result := range results.ExamResult {
		totalScore += result.Score
	}

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