package quiz_answer_repository

import (
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type quizAnswerRepository struct {
	database           *mongo.Database
	collectionQuestion string
	collectionAnswer   string
	collectionQuiz     string
}

func (q *quizAnswerRepository) FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (quiz_answer_domain.Response, error) {
	collectionAnswer := q.database.Collection(q.collectionAnswer)
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}

	filter := bson.M{"question_id": idQuestion, "user_id": idUser}
	cursor, err := collectionAnswer.Find(ctx, filter)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var answers []quiz_answer_domain.QuizAnswer
	for cursor.Next(ctx) {
		var answer quiz_answer_domain.QuizAnswer
		if err = cursor.Decode(&answer); err != nil {
			return quiz_answer_domain.Response{}, err
		}

		// Gắn CourseID vào bài học
		answer.QuestionID = idQuestion
		answers = append(answers, answer)
	}

	response := quiz_answer_domain.Response{
		QuizAnswer: answers,
	}

	return response, nil
}

func (q *quizAnswerRepository) DeleteAllAnswerByExamID(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (q *quizAnswerRepository) CreateOne(ctx context.Context, quizAnswer *quiz_answer_domain.QuizAnswer) error {
	collectionAnswer := q.database.Collection(q.collectionAnswer)
	collectionQuestion := q.database.Collection(q.collectionQuestion)

	filterQuestionID := bson.M{"question_id": quizAnswer.QuestionID}
	countQuestionID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}
	if countQuestionID == 0 {
		return errors.New("the question ID do not exist")
	}

	_, err = collectionAnswer.InsertOne(ctx, quizAnswer)
	return nil
}

func (q *quizAnswerRepository) DeleteOne(ctx context.Context, quizID string) error {
	collectionAnswer := q.database.Collection(q.collectionAnswer)
	objID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionAnswer.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`quiz answer is removed`)
	}

	_, err = collectionAnswer.DeleteOne(ctx, filter)
	return err
}

func NewQuizAnswerRepository(db *mongo.Database, collectionQuestion string, collectionAnswer string, collectionQuiz string) quiz_answer_domain.IQuizAnswerRepository {
	return &quizAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
		collectionQuiz:     collectionQuiz,
	}
}
