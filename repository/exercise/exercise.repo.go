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

func (e *exerciseRepository) FetchManyByLessonID(ctx context.Context, unitID string) ([]exercise_domain.ExerciseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (e *exerciseRepository) FetchManyByUnitID(ctx context.Context, unitID string) ([]exercise_domain.ExerciseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (e *exerciseRepository) UpdateCompleted(ctx context.Context, exerciseID string, isComplete int) error {
	//TODO implement me
	panic("implement me")
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
		var exercise exercise_domain.Exercise
		if err = cursor.Decode(&exercise); err != nil {
			return exercise_domain.ExerciseResponse{}, err
		}

		// Lấy thông tin liên quan cho mỗi khóa học
		countQuest, err := e.countQuestionByExerciseID(ctx, exercise.Id)
		if err != nil {
			return exercise_domain.ExerciseResponse{}, err
		}

		exerciseRes := exercise_domain.ExerciseResponse{
			ID:            exercise.Id,
			LessonID:      exercise.LessonID,
			UnitID:        exercise.UnitID,
			Title:         exercise.Title,
			Description:   exercise.Description,
			Duration:      exercise.Duration,
			CreatedAt:     exercise.CreatedAt,
			UpdatedAt:     exercise.UpdatedAt,
			WhoUpdates:    exercise.WhoUpdates,
			CountQuestion: countQuest,
		}

		exercises = append(exercises, exerciseRes)
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

func (e *exerciseRepository) FetchMany(ctx context.Context, page string) ([]exercise_domain.ExerciseResponse, int64, error) {
	collectionExercise := e.database.Collection(e.collectionExercise)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, errors.New("invalid page number")
	}
	perPage := 1
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	count, err := collectionExercise.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	cursor, err := collectionExercise.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var exercises []exercise_domain.ExerciseResponse

	for cursor.Next(ctx) {
		var exercise exercise_domain.Exercise
		if err := cursor.Decode(&exercise); err != nil {
			return nil, 0, err
		}

		// Lấy thông tin liên quan cho mỗi khóa học
		countQuest, err := e.countQuestionByExerciseID(ctx, exercise.Id)
		if err != nil {
			return nil, 0, err
		}

		exerciseRes := exercise_domain.ExerciseResponse{
			ID:            exercise.Id,
			LessonID:      exercise.LessonID,
			UnitID:        exercise.UnitID,
			Title:         exercise.Title,
			Description:   exercise.Description,
			Duration:      exercise.Duration,
			CreatedAt:     exercise.CreatedAt,
			UpdatedAt:     exercise.UpdatedAt,
			WhoUpdates:    exercise.WhoUpdates,
			CountQuestion: countQuest,
		}

		exercises = append(exercises, exerciseRes)
	}

	return exercises, cal, nil
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
