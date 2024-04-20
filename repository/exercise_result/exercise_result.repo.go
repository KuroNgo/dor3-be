package exercise_result

import (
	exercise_result_domain "clean-architecture/domain/exercise_result"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type exerciseResultRepository struct {
	database                 mongo.Database
	collectionExerciseResult string
	collectionExercise       string
}

func NewExerciseResultRepository(db mongo.Database, collectionExerciseResult string, collectionExercise string) exercise_result_domain.IExerciseResultRepository {
	return &exerciseResultRepository{
		database:                 db,
		collectionExerciseResult: collectionExerciseResult,
		collectionExercise:       collectionExercise,
	}
}

func (e *exerciseResultRepository) FetchMany(ctx context.Context, page string) (exercise_result_domain.Response, error) {
	collectionResult := e.database.Collection(e.collectionExerciseResult)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return exercise_result_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	cursor, err := collectionResult.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return exercise_result_domain.Response{}, err
	}

	var results []exercise_result_domain.ExerciseResult
	for cursor.Next(ctx) {
		var result exercise_result_domain.ExerciseResult
		if err = cursor.Decode(&result); err != nil {
			return exercise_result_domain.Response{}, err
		}

		results = append(results, result)
	}

	resultRes := exercise_result_domain.Response{
		ExerciseResult: results,
	}

	return resultRes, nil
}

func (e *exerciseResultRepository) FetchManyByExerciseID(ctx context.Context, exerciseID string) (exercise_result_domain.Response, error) {
	collectionResult := e.database.Collection(e.collectionExerciseResult)

	idExam, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return exercise_result_domain.Response{}, err
	}

	filter := bson.M{"exam_id": idExam}
	cursor, err := collectionResult.Find(ctx, filter)
	if err != nil {
		return exercise_result_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var results []exercise_result_domain.ExerciseResult
	for cursor.Next(ctx) {
		var result exercise_result_domain.ExerciseResult
		if err = cursor.Decode(&result); err != nil {
			return exercise_result_domain.Response{}, err
		}

		result.ExerciseID = idExam
		results = append(results, result)
	}

	resultRes := exercise_result_domain.Response{
		ExerciseResult: results,
	}
	return resultRes, nil
}

func (e *exerciseResultRepository) FetchManyByUserID(ctx context.Context, userID string) (exercise_result_domain.Response, error) {
	collectionResult := e.database.Collection(e.collectionExerciseResult)

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return exercise_result_domain.Response{}, err
	}

	filter := bson.M{"user_id": idUser}
	cursor, err := collectionResult.Find(ctx, filter)
	if err != nil {
		return exercise_result_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var results []exercise_result_domain.ExerciseResult
	score, _ := e.GetScoreByUser(ctx, userID)

	for cursor.Next(ctx) {
		var result exercise_result_domain.ExerciseResult
		if err = cursor.Decode(&result); err != nil {
			return exercise_result_domain.Response{}, err
		}

		result.UserID = idUser
		result.Score = score
		results = append(results, result)
	}
	averageScore, _ := e.GetAverageScoreByUser(ctx, userID)
	percentScore, _ := e.GetOverallPerformance(ctx, userID)

	resultRes := exercise_result_domain.Response{
		ExerciseResult: results,
		AverageScore:   averageScore,
		Percentage:     percentScore,
	}

	return resultRes, nil
}

func (e *exerciseResultRepository) GetResultsByUserIDAndExerciseID(ctx context.Context, userID string, exerciseID string) (exercise_result_domain.ExerciseResult, error) {
	//TODO implement me
	panic("implement me")
}

func (e *exerciseResultRepository) GetAverageScoreByUser(ctx context.Context, userID string) (float64, error) {
	results, err := e.FetchManyByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	var totalScore int16 = 0
	for _, result := range results.ExerciseResult {
		totalScore += result.Score
	}

	if len(results.ExerciseResult) > 0 {
		averageScore := float64(totalScore) / float64(len(results.ExerciseResult))
		return averageScore, nil
	}

	return 0, nil
}

func (e *exerciseResultRepository) GetOverallPerformance(ctx context.Context, userID string) (float64, error) {
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

func (e *exerciseResultRepository) CreateOne(ctx context.Context, exerciseResult *exercise_result_domain.ExerciseResult) error {
	collectionResult := e.database.Collection(e.collectionExerciseResult)
	collectionExam := e.database.Collection(e.collectionExercise)

	filterExamID := bson.M{"exercise_id": exerciseResult.ExerciseID}
	countLessonID, err := collectionExam.CountDocuments(ctx, filterExamID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the exercise ID do not exist")
	}

	_, err = collectionResult.InsertOne(ctx, exerciseResult)
	return nil
}

func (e *exerciseResultRepository) UpdateStatus(ctx context.Context, exerciseResultID string, status int) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionExerciseResult)

	filter := bson.D{{Key: "exam_id", Value: exerciseResultID}}
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

func (e *exerciseResultRepository) DeleteOne(ctx context.Context, exerciseResultID string) error {
	collectionResult := e.database.Collection(e.collectionExerciseResult)

	objID, err := primitive.ObjectIDFromHex(exerciseResultID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionResult.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exercise answer is removed`)
	}

	_, err = collectionResult.DeleteOne(ctx, filter)
	return err
}

func (e *exerciseResultRepository) CalculateScore(ctx context.Context, correctAnswers, totalQuestions int) int {
	if totalQuestions == 0 {
		return 0
	}

	score := (correctAnswers * 10) / totalQuestions
	return score
}

func (e *exerciseResultRepository) CalculatePercentage(ctx context.Context, correctAnswers, totalQuestions int) float64 {
	if totalQuestions == 0 {
		return 0
	}

	percentage := float64(correctAnswers) / float64(totalQuestions) * 100
	return percentage
}

func (e *exerciseResultRepository) GetScoreByUser(ctx context.Context, userID string) (int16, error) {
	results, err := e.FetchManyByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	var totalScore int16 = 0
	for _, result := range results.ExerciseResult {
		totalScore += result.Score
	}

	return totalScore, nil
}
