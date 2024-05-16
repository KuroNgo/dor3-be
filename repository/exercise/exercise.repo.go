package exercise_repository

import (
	exercise_domain "clean-architecture/domain/exercise"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
	collectionUnit := e.database.Collection(e.collectionUnit)
	collectionLesson := e.database.Collection(e.collectionLesson)

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
	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			var exercise exercise_domain.Exercise
			if err = cursor.Decode(&exercise); err != nil {
				return
			}

			// Fetch related data
			countQuest, err := e.countQuestionByExerciseID(ctx, exercise.Id)
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

			var exerciseRes exercise_domain.ExerciseResponse
			exerciseRes.ID = exercise.Id
			exerciseRes.Title = exercise.Title
			exerciseRes.Description = exercise.Description
			exerciseRes.Duration = exercise.Duration
			exerciseRes.CreatedAt = exercise.CreatedAt
			exerciseRes.UpdatedAt = exercise.UpdatedAt
			exerciseRes.WhoUpdates = exercise.WhoUpdates
			exerciseRes.Learner = exercise.Learner
			exerciseRes.IsComplete = exercise.IsComplete
			exerciseRes.CountQuestion = countQuest
			exerciseRes.Unit = unit
			exerciseRes.Lesson = lesson

			exercises = append(exercises, exerciseRes)
		}
	}()
	internal.Wg.Wait()

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
	collectionUnit := e.database.Collection(e.collectionUnit)
	collectionLesson := e.database.Collection(e.collectionLesson)

	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		return nil, exercise_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.D{{"_id", -1}})

	// Convert unitID to ObjectID
	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	// Count documents for pagination
	count, err := collectionExercise.CountDocuments(ctx, bson.M{"unit_id": idUnit})
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	totalPages := (count + int64(perPage) - 1) / int64(perPage) // Calculate total pages

	// Query for exercises
	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionExercise.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}
	defer cursor.Close(ctx)

	var exercises []exercise_domain.ExerciseResponse

	// Process each exercise
	for cursor.Next(ctx) {
		var exercise exercise_domain.Exercise
		if err = cursor.Decode(&exercise); err != nil {
			return nil, exercise_domain.DetailResponse{}, err
		}

		// Fetch related data
		countQuest, err := e.countQuestionByExerciseID(ctx, exercise.Id)
		if err != nil {
			return nil, exercise_domain.DetailResponse{}, err
		}

		var unit unit_domain.Unit
		if err = collectionUnit.FindOne(ctx, bson.M{"_id": idUnit}).Decode(&unit); err != nil {
			return nil, exercise_domain.DetailResponse{}, err
		}

		var lesson lesson_domain.Lesson
		if err = collectionLesson.FindOne(ctx, bson.M{"_id": unit.LessonID}).Decode(&lesson); err != nil {
			return nil, exercise_domain.DetailResponse{}, err
		}

		var exerciseRes exercise_domain.ExerciseResponse
		exerciseRes.ID = exercise.Id
		exerciseRes.Title = exercise.Title
		exerciseRes.Description = exercise.Description
		exerciseRes.Duration = exercise.Duration
		exerciseRes.CreatedAt = exercise.CreatedAt
		exerciseRes.UpdatedAt = exercise.UpdatedAt
		exerciseRes.WhoUpdates = exercise.WhoUpdates
		exerciseRes.Learner = exercise.Learner
		exerciseRes.IsComplete = exercise.IsComplete
		exerciseRes.CountQuestion = countQuest
		exerciseRes.Unit = unit
		exerciseRes.Lesson = lesson

		exercises = append(exercises, exerciseRes)
	}

	if err = cursor.Err(); err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	detail := exercise_domain.DetailResponse{
		CountExercise: count,
		Page:          totalPages,
		CurrentPage:   pageNumber,
	}

	return exercises, detail, nil
}

func (e *exerciseRepository) FetchMany(ctx context.Context, page string) ([]exercise_domain.ExerciseResponse, exercise_domain.DetailResponse, error) {
	collectionExercise := e.database.Collection(e.collectionExercise)
	collectionUnit := e.database.Collection(e.collectionUnit)
	collectionLesson := e.database.Collection(e.collectionLesson)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 1
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))
	count, err := collectionExercise.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}

	// Calculate total pages directly without goroutine
	totalPages := (count + int64(perPage) - 1) / int64(perPage)

	cursor, err := collectionExercise.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, exercise_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("failed to close cursor: %v", err)
		}
	}(cursor, ctx)

	var exercises []exercise_domain.ExerciseResponse

	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			var exercise exercise_domain.Exercise
			if err := cursor.Decode(&exercise); err != nil {
				log.Printf("failed to decode exercise: %v", err)
				return
			}

			// Lấy thông tin liên quan cho mỗi khóa học
			countQuest, err := e.countQuestionByExerciseID(ctx, exercise.Id)
			if err != nil {
				log.Printf("failed to count questions: %v", err)
				return
			}

			var unit unit_domain.Unit
			filterUnit := bson.M{"_id": exercise.UnitID}
			err = collectionUnit.FindOne(ctx, filterUnit).Decode(&unit)
			if err != nil {
				log.Printf("failed to find unit: %v", err)
				return
			}

			var lesson lesson_domain.Lesson
			filterLesson := bson.M{"_id": unit.LessonID}
			err = collectionLesson.FindOne(ctx, filterLesson).Decode(&lesson)
			if err != nil {
				log.Printf("failed to find lesson: %v", err)
				return
			}

			var exerciseRes exercise_domain.ExerciseResponse
			exerciseRes.ID = exercise.Id
			exerciseRes.Title = exercise.Title
			exerciseRes.Description = exercise.Description
			exerciseRes.Duration = exercise.Duration
			exerciseRes.CreatedAt = exercise.CreatedAt
			exerciseRes.UpdatedAt = exercise.UpdatedAt
			exerciseRes.WhoUpdates = exercise.WhoUpdates
			exerciseRes.Learner = exercise.Learner
			exerciseRes.IsComplete = exercise.IsComplete
			exerciseRes.CountQuestion = countQuest
			exerciseRes.Unit = unit
			exerciseRes.Lesson = lesson

			exercises = append(exercises, exerciseRes)
		}
		if err := cursor.Err(); err != nil {
			log.Printf("cursor error: %v", err)
		}
	}()
	internal.Wg.Wait()

	statisticsCh := make(chan exercise_domain.Statistics)
	go func() {
		statistics, _ := e.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	detail := exercise_domain.DetailResponse{
		CountExercise: count,
		Page:          totalPages,
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
	collectionLesson := e.database.Collection(e.collectionLesson)
	collectionUnit := e.database.Collection(e.collectionUnit)

	filterLessonID := bson.M{"_id": exercise.LessonID}
	countLessonID, err := collectionLesson.CountDocuments(ctx, filterLessonID)
	if err != nil {
		return err
	}
	if countLessonID == 0 {
		return errors.New("the lesson ID does not exist")
	}

	filterUnitID := bson.M{"_id": exercise.UnitID}
	countUnitID, err := collectionUnit.CountDocuments(ctx, filterUnitID)
	if err != nil {
		return err
	}
	if countUnitID == 0 {
		return errors.New("the unit ID does not exist")
	}

	_, err = collectionExercise.InsertOne(ctx, exercise)
	if err != nil {
		return err
	}

	return nil
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
