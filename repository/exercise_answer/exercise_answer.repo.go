package exercise_answer_repository

import (
	"clean-architecture/domain/exercise_answer"
	"clean-architecture/internal"
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
	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			var answer exercise_answer_domain.ExerciseAnswer
			if err = cursor.Decode(&answer); err != nil {
				return
			}

			answers = append(answers, answer)
		}
	}()

	internal.Wg.Wait()

	response := exercise_answer_domain.Response{
		ExerciseAnswer: answers,
	}

	return response, nil
}

func (e *exerciseAnswerRepository) CreateOne(ctx context.Context, exerciseAnswer *exercise_answer_domain.ExerciseAnswer) error {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	// kiểm tra questionId có tồn tại
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

func (e *exerciseAnswerRepository) DeleteAllAnswerByExerciseID(ctx context.Context, exerciseId string) error {
	collectionAnswer := e.database.Collection(e.collectionAnswer)

	objID, err := primitive.ObjectIDFromHex(exerciseId)
	if err != nil {
		return err
	}

	filter := bson.M{"exercise_id": objID}
	count, err := collectionAnswer.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exercise answer is removed`)
	}

	_, err = collectionAnswer.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
