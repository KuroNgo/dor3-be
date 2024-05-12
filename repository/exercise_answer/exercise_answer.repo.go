package exercise_answer_repository

import (
	"clean-architecture/domain/exercise_answer"
	exercise_options_domain "clean-architecture/domain/exercise_options"
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
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
	collectionOptions  string
}

func NewExerciseAnswerRepository(db *mongo.Database, collectionQuestion string, collectionAnswer string, collectionOptions string) exercise_answer_domain.IExerciseAnswerRepository {
	return &exerciseAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
		collectionOptions:  collectionOptions,
	}
}

func (e *exerciseAnswerRepository) FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (exercise_answer_domain.Response, error) {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectionOptions := e.database.Collection(e.collectionOptions)

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

	var answers []exercise_answer_domain.ExerciseAnswerResponse
	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			var answer exercise_answer_domain.ExerciseAnswer
			if err = cursor.Decode(&answer); err != nil {
				return
			}

			var question exercise_questions_domain.ExerciseQuestion
			filterQuestion := bson.M{"_id": answer.QuestionID}
			err = collectionQuestion.FindOne(ctx, filterQuestion).Decode(&question)
			if err != nil {
				return
			}

			var options exercise_options_domain.ExerciseOptions
			filterOptions := bson.M{"question_id": question.ID}
			err = collectionOptions.FindOne(ctx, filterOptions).Decode(&options)
			if err != nil {
				return
			}

			var answerRes exercise_answer_domain.ExerciseAnswerResponse
			if err = cursor.Decode(&answerRes); err != nil {
				return
			}

			answerRes.ID = answer.ID
			answerRes.UserID = answer.UserID
			answerRes.Question = question
			answerRes.Options = options
			answerRes.IsCorrect = answer.IsCorrect
			answerRes.BlankIndex = answer.BlankIndex
			answerRes.SubmittedAt = answer.SubmittedAt

			answers = append(answers, answerRes)
		}
	}()

	internal.Wg.Wait()

	response := exercise_answer_domain.Response{
		ExerciseAnswerResponse: answers,
	}

	return response, nil
}

func (e *exerciseAnswerRepository) CreateOne(ctx context.Context, exerciseAnswer *exercise_answer_domain.ExerciseAnswer) error {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionOptions := e.database.Collection(e.collectionOptions)
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

	// kiểm tra answer có bằng với đáp án
	var options exercise_options_domain.ExerciseOptions
	if err := collectionOptions.FindOne(ctx, filterQuestionID).Decode(&options); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("no options found for the question ID")
		}
		return err
	}

	if exerciseAnswer.Answer == options.CorrectAnswer {
		exerciseAnswer.IsCorrect = 1 //đúng
	} else {
		exerciseAnswer.IsCorrect = 0 //sai
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
