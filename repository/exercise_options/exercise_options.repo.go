package exercise_options_repository

import (
	exercise_options_domain "clean-architecture/domain/exercise_options"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type exerciseOptionsRepository struct {
	database           *mongo.Database
	collectionQuestion string
	collectionOptions  string
}

func (e *exerciseOptionsRepository) FetchManyByQuestionID(ctx context.Context, questionID string) (exercise_options_domain.Response, error) {
	collectionOptions := e.database.Collection(e.collectionOptions)
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return exercise_options_domain.Response{}, err
	}

	filter := bson.M{"question_id": idQuestion}
	cursor, err := collectionOptions.Find(ctx, filter)
	if err != nil {
		return exercise_options_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var options []exercise_options_domain.ExerciseOptions
	for cursor.Next(ctx) {
		var option exercise_options_domain.ExerciseOptions
		if err = cursor.Decode(&option); err != nil {
			return exercise_options_domain.Response{}, err
		}

		// Gắn CourseID vào bài học
		option.QuestionID = idQuestion
		options = append(options, option)
	}

	response := exercise_options_domain.Response{
		ExerciseOptions: options,
	}

	return response, nil
}

func (e *exerciseOptionsRepository) UpdateOne(ctx context.Context, exerciseOptions *exercise_options_domain.ExerciseOptions) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionQuestion)

	filter := bson.D{{Key: "_id", Value: exerciseOptions.ID}}
	update := bson.M{
		"$set": bson.M{
			"question_id": exerciseOptions.QuestionID,
			"content":     exerciseOptions.Content,
			"update_at":   exerciseOptions.UpdateAt,
			"who_update":  exerciseOptions.WhoUpdate,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *exerciseOptionsRepository) CreateOne(ctx context.Context, exerciseOptions *exercise_options_domain.ExerciseOptions) error {
	collectionOptions := e.database.Collection(e.collectionOptions)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	filterQuestionID := bson.M{"question_id": exerciseOptions.QuestionID}
	countLessonID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the question ID do not exist")
	}

	_, err = collectionOptions.InsertOne(ctx, exerciseOptions)
	return nil
}

func (e *exerciseOptionsRepository) DeleteOne(ctx context.Context, optionsID string) error {
	collectionOptions := e.database.Collection(e.collectionOptions)
	objID, err := primitive.ObjectIDFromHex(optionsID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionOptions.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exercise answer is removed`)
	}

	_, err = collectionOptions.DeleteOne(ctx, filter)
	return err
}

func NewExamOptionsRepository(db *mongo.Database, collectionQuestion string, collectionOptions string) exercise_options_domain.IExerciseOptionRepository {
	return &exerciseOptionsRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionOptions:  collectionOptions,
	}
}
