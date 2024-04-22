package quiz_result_repository

import (
	quiz_result_domain "clean-architecture/domain/quiz_result"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

type quizResultRepository struct {
	database             *mongo.Database
	collectionQuizResult string
	collectionQuiz       string
}

func NewQuizResultRepository(db *mongo.Database, collectionQuizResult string, collectionQuiz string) quiz_result_domain.IQuizResultRepository {
	return &quizResultRepository{
		database:             db,
		collectionQuizResult: collectionQuizResult,
		collectionQuiz:       collectionQuiz,
	}
}

func (q *quizResultRepository) FetchMany(ctx context.Context, page string) (quiz_result_domain.Response, error) {
	collectionResult := q.database.Collection(q.collectionQuizResult)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return quiz_result_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	cursor, err := collectionResult.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return quiz_result_domain.Response{}, err
	}

	var results []quiz_result_domain.QuizResult
	for cursor.Next(ctx) {
		var result quiz_result_domain.QuizResult
		if err = cursor.Decode(&result); err != nil {
			return quiz_result_domain.Response{}, err
		}

		results = append(results, result)
	}

	resultRes := quiz_result_domain.Response{
		QuizResult: results,
	}

	return resultRes, nil
}

func (q *quizResultRepository) FetchManyByQuizID(ctx context.Context, quizID string) (quiz_result_domain.Response, error) {
	collectionResult := q.database.Collection(q.collectionQuizResult)

	idQuiz, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return quiz_result_domain.Response{}, err
	}

	filter := bson.M{"quiz_id": idQuiz}
	cursor, err := collectionResult.Find(ctx, filter)
	if err != nil {
		return quiz_result_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var results []quiz_result_domain.QuizResult
	for cursor.Next(ctx) {
		var result quiz_result_domain.QuizResult
		if err = cursor.Decode(&result); err != nil {
			return quiz_result_domain.Response{}, err
		}

		result.QuizID = idQuiz
		results = append(results, result)
	}

	resultRes := quiz_result_domain.Response{
		QuizResult: results,
	}
	return resultRes, nil
}

func (q *quizResultRepository) FetchManyByUserID(ctx context.Context, userID string) (quiz_result_domain.Response, error) {
	collectionResult := q.database.Collection(q.collectionQuizResult)

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return quiz_result_domain.Response{}, err
	}

	filter := bson.M{"user_id": idUser}
	cursor, err := collectionResult.Find(ctx, filter)
	if err != nil {
		return quiz_result_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var results []quiz_result_domain.QuizResult
	score, _ := q.GetScoreByUser(ctx, userID)

	for cursor.Next(ctx) {
		var result quiz_result_domain.QuizResult
		if err = cursor.Decode(&result); err != nil {
			return quiz_result_domain.Response{}, err
		}

		result.UserID = idUser
		result.Score = score
		results = append(results, result)
	}
	averageScore, _ := q.GetAverageScoreByUser(ctx, userID)
	percentScore, _ := q.GetOverallPerformance(ctx, userID)

	resultRes := quiz_result_domain.Response{
		QuizResult:   results,
		AverageScore: averageScore,
		Percentage:   percentScore,
	}

	return resultRes, nil
}

func (q *quizResultRepository) GetResultsByUserIDAndQuizID(ctx context.Context, userID string, exerciseID string) (quiz_result_domain.QuizResult, error) {
	collection := q.database.Collection(q.collectionQuizResult)

	var result quiz_result_domain.QuizResult

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return quiz_result_domain.QuizResult{}, err
	}

	idExam, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return quiz_result_domain.QuizResult{}, err
	}

	filter := bson.M{"user_id": idUser, "exam_id": idExam}

	err = collection.FindOne(ctx, filter).Decode(&result)
	return result, err
}

func (q *quizResultRepository) GetAverageScoreByUser(ctx context.Context, userID string) (float64, error) {
	results, err := q.FetchManyByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	var totalScore int16 = 0
	for _, result := range results.QuizResult {
		totalScore += result.Score
	}

	if len(results.QuizResult) > 0 {
		averageScore := float64(totalScore) / float64(len(results.QuizResult))
		return averageScore, nil
	}

	return 0, nil
}

func (e *quizResultRepository) GetScoreByUser(ctx context.Context, userID string) (int16, error) {
	// Lấy tất cả các kết quả của bài kiểm tra từ cơ sở dữ liệu MongoDB
	results, err := e.FetchManyByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Tính tổng điểm
	var totalScore int16 = 0
	for _, result := range results.QuizResult {
		totalScore += result.Score
	}

	return totalScore, nil
}

func (q *quizResultRepository) GetOverallPerformance(ctx context.Context, userID string) (float64, error) {
	averageScore, err := q.GetAverageScoreByUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	if averageScore == 0 {
		return 0, nil
	}

	overallPerformance := averageScore * 100

	return overallPerformance, nil
}

func (q *quizResultRepository) CreateOne(ctx context.Context, quizResult *quiz_result_domain.QuizResult) error {
	collectionResult := q.database.Collection(q.collectionQuizResult)
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	filterQuizID := bson.M{"quiz_id": quizResult.QuizID}
	countLessonID, err := collectionQuiz.CountDocuments(ctx, filterQuizID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the quizID do not exist")
	}

	_, err = collectionResult.InsertOne(ctx, quizResult)
	return nil
}

func (q *quizResultRepository) DeleteOne(ctx context.Context, quizResultID string) error {
	collectionResult := q.database.Collection(q.collectionQuizResult)

	objID, err := primitive.ObjectIDFromHex(quizResultID)
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

func (q *quizResultRepository) UpdateStatus(ctx context.Context, quizResultID string, status int) (*mongo.UpdateResult, error) {
	collection := q.database.Collection(q.collectionQuizResult)

	filter := bson.D{{Key: "quiz_id", Value: quizResultID}}
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

func (q *quizResultRepository) CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int {
	if totalQuestions == 0 {
		return 0
	}

	score := (correctAnswers * 10) / totalQuestions
	return score
}

func (q *quizResultRepository) CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64 {
	if totalQuestions == 0 {
		return 0
	}

	percentage := float64(correctAnswers) / float64(totalQuestions) * 100
	return percentage
}
