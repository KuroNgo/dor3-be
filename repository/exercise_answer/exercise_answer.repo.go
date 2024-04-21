package exercise_answer_repository

import (
	"clean-architecture/domain/exercise_answer"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type exerciseAnswerRepository struct {
	database           *mongo.Database
	collectionQuestion string
	collectionAnswer   string
}

func NewExerciseAnswerRepository(db *mongo.Database, collectionQuestion string, collectionAnswer string) exercise_answer_domain.IExerciseAnswerRepository {
	return &exerciseAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
	}
}

func (e *exerciseAnswerRepository) FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (exercise_answer_domain.Response, error) {
	collectionAnswer := e.database.Collection(e.collectionAnswer)

	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return exercise_answer_domain.Response{}, err
	}

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return exercise_answer_domain.Response{}, err
	}

	filter := bson.M{"question_id": idQuestion, "user_id": idUser}
	cursor, err := collectionAnswer.Find(ctx, filter)
	if err != nil {
		return exercise_answer_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var answers []exercise_answer_domain.ExerciseAnswer
	for cursor.Next(ctx) {
		var answer exercise_answer_domain.ExerciseAnswer
		if err = cursor.Decode(&answer); err != nil {
			return exercise_answer_domain.Response{}, err
		}

		// Gắn CourseID vào bài học
		answer.QuestionID = idQuestion
		answers = append(answers, answer)
	}

	response := exercise_answer_domain.Response{
		ExerciseAnswer: answers,
	}

	return response, nil
}

func (e *exerciseAnswerRepository) CreateOne(ctx context.Context, exerciseAnswer *exercise_answer_domain.ExerciseAnswer) error {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	filterQuestionID := bson.M{"question_id": exerciseAnswer.QuestionID}
	countLessonID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the question ID do not exist")
	}

	_, err = collectionAnswer.InsertOne(ctx, exerciseAnswer)
	return nil
}

func (e *exerciseAnswerRepository) DeleteOne(ctx context.Context, exerciseID string) error {
	collectionExercise := e.database.Collection(e.collectionAnswer)
	objID, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionExercise.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam answer is removed`)
	}

	_, err = collectionExercise.DeleteOne(ctx, filter)
	return err
}
