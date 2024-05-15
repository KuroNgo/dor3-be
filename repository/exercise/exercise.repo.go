package exercise_repository

import (
	exercise_domain "clean-architecture/domain/exercise"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"strconv"
	"time"
)

type exerciseRepository struct {
	database             *mongo.Database
	collectionLesson     string
	collectionUnit       string
	collectionVocabulary string
	collectionExercise   string
	collectionQuestion   string
}

func NewExerciseRepository(db *mongo.Database, collectionLesson string, collectionUnit string, collectionVocabulary string, collectionExercise string, collectionQuestion string) exercise_domain.IExerciseRepository {
	return &exerciseRepository{
		database:             db,
		collectionLesson:     collectionLesson,
		collectionUnit:       collectionUnit,
		collectionVocabulary: collectionVocabulary,
		collectionExercise:   collectionExercise,
		collectionQuestion:   collectionQuestion,
	}
}

func (e *exerciseRepository) FetchOneByUnitID(ctx context.Context, unitID string) (exercise_domain.ExerciseResponse, error) {
	collectionExercise := e.database.Collection(e.collectionExercise)

	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return exercise_domain.ExerciseResponse{}, err
	}

	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionExercise.Find(ctx, filter)
	if err != nil {
		return exercise_domain.ExerciseResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var exercises []exercise_domain.ExerciseResponse
	for cursor.Next(ctx) {
		var exercise exercise_domain.ExerciseResponse
		if err = cursor.Decode(&exercise); err != nil {
			return exercise_domain.ExerciseResponse{}, err
		}

		// Lấy thông tin liên quan cho mỗi khóa học
		countQuest, err := e.countQuestionByExerciseID(ctx, exercise.ID)
		if err != nil {
			return exercise_domain.ExerciseResponse{}, err
		}

		exercise.CountQuestion = countQuest

		exercises = append(exercises, exercise)
	}

	// Kiểm tra nếu danh sách exercises không rỗng
	if len(exercises) == 0 {
		return exercise_domain.ExerciseResponse{}, errors.New("no exercises found")
	}

	// Chọn một giá trị ngẫu nhiên từ danh sách exercises
	randomIndex := rand.Intn(len(exercises))
	randomExercise := exercises[randomIndex]

	return randomExercise, nil
}

func (e *exerciseRepository) FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]exercise_domain.ExerciseResponse, exercise_domain.DetailResponse, error) {
	collectionExercise := e.database.Collection(e.collectionExercise)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.D{{"_id", -1}})

	calCh := make(chan int64)
	countCh := make(chan int64)
	go func() {
		count, err := collectionExercise.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}
		countCh <- count

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1
		}
	}()

	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionExercise.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var exercises []exercise_domain.ExerciseResponse
	for cursor.Next(ctx) {
		var exercise exercise_domain.ExerciseResponse
		if err = cursor.Decode(&exercise); err != nil {
			return nil, exercise_domain.DetailResponse{}, err
		}

		// Lấy thông tin liên quan cho mỗi khóa học
		countQuest, err := e.countQuestionByExerciseID(ctx, exercise.ID)
		if err != nil {
			return nil, exercise_domain.DetailResponse{}, err
		}

		exercise.CountQuestion = countQuest

		exercises = append(exercises, exercise)
	}

	cal := <-calCh
	countExercise := <-countCh

	detail := exercise_domain.DetailResponse{
		CountExercise: countExercise,
		Page:          cal,
		CurrentPage:   pageNumber,
	}

	return exercises, detail, nil
}

func (e *exerciseRepository) FetchMany(ctx context.Context, page string) ([]exercise_domain.ExerciseResponse, exercise_domain.DetailResponse, error) {
	collectionExercise := e.database.Collection(e.collectionExercise)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 1
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	calCh := make(chan int64)
	countExerciseCh := make(chan int64)

	go func() {
		defer close(calCh)
		defer close(countExerciseCh)

		count, err := collectionExercise.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}

		countExerciseCh <- count

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1 + 1
		}
	}()

	cursor, err := collectionExercise.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var exercises []exercise_domain.ExerciseResponse

	for cursor.Next(ctx) {
		var exercise exercise_domain.ExerciseResponse
		if err := cursor.Decode(&exercise); err != nil {
			return nil, exercise_domain.DetailResponse{}, err
		}

		// Lấy thông tin liên quan cho mỗi khóa học
		countQuest, err := e.countQuestionByExerciseID(ctx, exercise.ID)
		if err != nil {
			return nil, exercise_domain.DetailResponse{}, err
		}

		exercise.CountQuestion = countQuest

		exercises = append(exercises, exercise)
	}

	cal := <-calCh
	countExercise := <-countExerciseCh
	statisticsCh := make(chan exercise_domain.Statistics)
	go func() {
		statistics, _ := e.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	detail := exercise_domain.DetailResponse{
		CountExercise: countExercise,
		Page:          cal,
		CurrentPage:   pageNumber,
		Statistics:    statistics,
	}

	return exercises, detail, nil
}

func (e *exerciseRepository) UpdateOne(ctx context.Context, exercise *exercise_domain.Exercise) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionExercise)

	filter := bson.D{{Key: "_id", Value: exercise.Id}}
	update := bson.M{
		"$set": bson.M{
			"lesson_id":  exercise.LessonID,
			"unit_id":    exercise.UnitID,
			"title":      exercise.Title,
			"duration":   exercise.Duration,
			"update_at":  exercise.UpdatedAt,
			"who_update": exercise.WhoUpdates,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *exerciseRepository) UpdateCompleted(ctx context.Context, exercise *exercise_domain.Exercise) error {
	collection := e.database.Collection(e.collectionExercise)

	filter := bson.D{{Key: "_id", Value: exercise.Id}}
	update := bson.M{
		"$set": bson.M{
			"is_complete": exercise.IsComplete,
			"update_at":   time.Now(),
			"learner":     exercise.Learner,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (e *exerciseRepository) CreateOne(ctx context.Context, exercise *exercise_domain.Exercise) error {
	collectionExercise := e.database.Collection(e.collectionExercise)
	_, err := collectionExercise.InsertOne(ctx, exercise)
	return err
}

func (e *exerciseRepository) DeleteOne(ctx context.Context, exerciseID string) error {
	collectionExercise := e.database.Collection(e.collectionExercise)

	objID, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}

	count, err := collectionExercise.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the exercise is removed or have not exist`)
	}

	_, err = collectionExercise.DeleteOne(ctx, filter)
	return err
}

// countLessonsByCourseID counts the number of lessons associated with a course.
func (e *exerciseRepository) countQuestionByExerciseID(ctx context.Context, exerciseId primitive.ObjectID) (int32, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	filter := bson.M{"exercise_id": exerciseId}
	count, err := collectionQuestion.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

func (e *exerciseRepository) Statistics(ctx context.Context) (exercise_domain.Statistics, error) {
	collectionExercise := e.database.Collection(e.collectionExercise)

	count, err := collectionExercise.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exercise_domain.Statistics{}, err
	}

	statistics := exercise_domain.Statistics{
		Total: count,
	}
	return statistics, nil
}
