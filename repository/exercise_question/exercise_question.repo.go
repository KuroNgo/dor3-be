package exercise_question_repository

import (
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type exerciseQuestionRepository struct {
	database           *mongo.Database
	collectionQuestion string
	collectionExercise string
}

func NewExerciseQuestionRepository(db *mongo.Database, collectionQuestion string, collectionExercise string) exercise_questions_domain.IExerciseQuestionRepository {
	return &exerciseQuestionRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionExercise: collectionExercise,
	}
}

func (e *exerciseQuestionRepository) FetchMany(ctx context.Context, page string) (exercise_questions_domain.Response, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return exercise_questions_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	count, err := collectionQuestion.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exercise_questions_domain.Response{}, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	cursor, err := collectionQuestion.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return exercise_questions_domain.Response{}, err
	}

	var questions []exercise_questions_domain.ExerciseQuestion
	for cursor.Next(ctx) {
		var question exercise_questions_domain.ExerciseQuestion
		if err = cursor.Decode(&question); err != nil {
			return exercise_questions_domain.Response{}, err
		}

		questions = append(questions, question)
	}
	questionsRes := exercise_questions_domain.Response{
		Page:             cal,
		ExerciseQuestion: questions,
	}
	return questionsRes, nil
}

func (e *exerciseQuestionRepository) FetchManyByExerciseID(ctx context.Context, exerciseID string) (exercise_questions_domain.Response, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	idExam, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return exercise_questions_domain.Response{}, err
	}

	filter := bson.M{"exercise_id": idExam}
	cursor, err := collectionQuestion.Find(ctx, filter)
	if err != nil {
		return exercise_questions_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var questions []exercise_questions_domain.ExerciseQuestion
	for cursor.Next(ctx) {
		var question exercise_questions_domain.ExerciseQuestion
		if err = cursor.Decode(&question); err != nil {
			return exercise_questions_domain.Response{}, err
		}

		question.ExerciseID = idExam
		questions = append(questions, question)
	}

	questionsRes := exercise_questions_domain.Response{
		ExerciseQuestion: questions,
	}

	return questionsRes, nil
}

func (e *exerciseQuestionRepository) UpdateOne(ctx context.Context, exerciseQuestion *exercise_questions_domain.ExerciseQuestion) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionQuestion)

	filter := bson.D{{Key: "_id", Value: exerciseQuestion.ID}}
	update := bson.M{
		"$set": bson.M{
			"exercise_id": exerciseQuestion.ExerciseID,
			"content":     exerciseQuestion.Content,
			"level":       exerciseQuestion.Level,
			"update_at":   exerciseQuestion.UpdateAt,
			"who_update":  exerciseQuestion.WhoUpdate,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *exerciseQuestionRepository) CreateOne(ctx context.Context, exerciseQuestion *exercise_questions_domain.ExerciseQuestion) error {
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectionExercise := e.database.Collection(e.collectionExercise)

	filterExerciseID := bson.M{"exercise_id": exerciseQuestion.ExerciseID}
	countExerciseID, err := collectionExercise.CountDocuments(ctx, filterExerciseID)
	if err != nil {
		return err
	}

	if countExerciseID == 0 {
		return errors.New("the exerciseID do not exist")
	}

	_, err = collectionQuestion.InsertOne(ctx, exerciseQuestion)
	return nil
}

func (e *exerciseQuestionRepository) DeleteOne(ctx context.Context, exerciseID string) error {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	objID, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionQuestion.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exercise answer is removed`)
	}

	_, err = collectionQuestion.DeleteOne(ctx, filter)
	return err
}
