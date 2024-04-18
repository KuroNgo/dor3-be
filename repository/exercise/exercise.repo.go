package exercise_repository

import (
	exercise_domain "clean-architecture/domain/exercise"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type exerciseRepository struct {
	database             *mongo.Database
	collectionLesson     string
	collectionUnit       string
	collectionVocabulary string
	collectionExercise   string
}

func NewExerciseRepository(db *mongo.Database, collectionVocabulary string, collectionExercise string) exercise_domain.IExerciseRepository {
	return &exerciseRepository{
		database:             db,
		collectionVocabulary: collectionVocabulary,
		collectionExercise:   collectionExercise,
	}
}

func (e *exerciseRepository) FetchManyByLessonID(ctx context.Context, unitID string) (exercise_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (e *exerciseRepository) UpdateCompleted(ctx context.Context, exerciseID string, isComplete int) error {
	//TODO implement me
	panic("implement me")
}

func (e *exerciseRepository) FetchManyByUnitID(ctx context.Context, unitID string) (exercise_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (e *exerciseRepository) FetchMany(ctx context.Context, page string) (exercise_domain.Response, error) {
	collectionExercise := e.database.Collection(e.collectionExercise)
	collectionVocabulary := e.database.Collection(e.collectionVocabulary)

	// Đếm tổng số lượng tài liệu trong collection
	count, err := collectionExercise.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exercise_domain.Response{}, err
	}

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return exercise_domain.Response{}, errors.New("invalid page number")
	}

	perPage := 1
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))
	cursor, err := collectionExercise.Find(ctx, bson.D{}, findOptions)

	if err != nil {
		return exercise_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var exercises []exercise_domain.Exercise

	for cursor.Next(ctx) {
		var exercise exercise_domain.Exercise

		if err := cursor.Decode(&exercise); err != nil {
			return exercise_domain.Response{}, err
		}

		var vocabulary vocabulary_domain.Vocabulary
		if err := collectionVocabulary.FindOne(ctx, bson.M{"_id": vocabulary.UnitID}).Decode(&vocabulary); err != nil {
			return exercise_domain.Response{}, err
		}

		exercise.VocabularyID = vocabulary.Id
		exercises = append(exercises, exercise)
	}

	exerciseRes := exercise_domain.Response{
		Exercise: exercises,
		Count:    count,
	}

	return exerciseRes, nil
}

func (e *exerciseRepository) UpdateOne(ctx context.Context, exercise *exercise_domain.Exercise) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionExercise)

	filter := bson.D{{Key: "_id", Value: exercise.Id}}
	update := bson.M{
		"$set": bson.M{
			"lesson_id":  exercise.LessonID,
			"unit_id":    exercise.UnitID,
			"vocabulary": exercise.VocabularyID,
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
	collectionVocabulary := e.database.Collection(e.collectionVocabulary)

	filter := bson.M{"content": exercise.Title, "vocabulary_id": exercise.VocabularyID}
	filterReference := bson.M{"_id": exercise.VocabularyID}

	countParent, err := collectionVocabulary.CountDocuments(ctx, filterReference)
	if err != nil {
		return err
	}

	// check exists with CountDocuments
	count, err := collectionExercise.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the content of exercise in vocabulary did exist")
	}
	if countParent == 0 {
		return errors.New("the vocabulary ID do not exist")
	}

	_, err = collectionVocabulary.InsertOne(ctx, exercise)
	return nil
}

func (e *exerciseRepository) UpsertOne(ctx context.Context, id string, exercise *exercise_domain.Exercise) (exercise_domain.Response, error) {
	collectionExercise := e.database.Collection(e.collectionExercise)
	collectionVocabulary := e.database.Collection(e.collectionVocabulary)

	filterReference := bson.M{"_id": exercise.VocabularyID}
	count, err := collectionVocabulary.CountDocuments(ctx, filterReference)
	if err != nil {
		return exercise_domain.Response{}, err
	}

	if count == 0 {
		return exercise_domain.Response{}, errors.New("the vocabulary ID do not exist")
	}
	doc, err := internal.ToDoc(exercise)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return exercise_domain.Response{}, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionExercise.FindOneAndUpdate(ctx, query, update, opts)

	var updatePost exercise_domain.Response
	if err := res.Decode(&updatePost); err != nil {
		return exercise_domain.Response{}, err
	}

	return updatePost, nil
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
