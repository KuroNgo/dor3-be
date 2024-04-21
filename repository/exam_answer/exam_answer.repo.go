package exam_answer_repository

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type examAnswerRepository struct {
	database           *mongo.Database
	collectionQuestion string
	collectionAnswer   string
	collectionExam     string
}

func NewExamAnswerRepository(db *mongo.Database, collectionQuestion string, collectionAnswer string, collectionExam string) exam_answer_domain.IExamAnswerRepository {
	return &examAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
		collectionExam:     collectionExam,
	}
}

func (e *examAnswerRepository) FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (exam_answer_domain.Response, error) {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	filter := bson.M{"question_id": idQuestion, "user_id": idUser}
	cursor, err := collectionAnswer.Find(ctx, filter)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var answers []exam_answer_domain.ExamAnswer
	for cursor.Next(ctx) {
		var answer exam_answer_domain.ExamAnswer
		if err = cursor.Decode(&answer); err != nil {
			return exam_answer_domain.Response{}, err
		}

		// Gắn CourseID vào bài học
		answer.QuestionID = idQuestion
		answers = append(answers, answer)
	}

	response := exam_answer_domain.Response{
		ExamAnswer: answers,
	}

	return response, nil
}

func (e *examAnswerRepository) CreateOne(ctx context.Context, examAnswer *exam_answer_domain.ExamAnswer) error {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	filterQuestionID := bson.M{"question_id": examAnswer.QuestionID}
	countQuestionID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}
	if countQuestionID == 0 {
		return errors.New("the question ID do not exist")
	}

	_, err = collectionAnswer.InsertOne(ctx, examAnswer)
	return nil
}

func (e *examAnswerRepository) DeleteOne(ctx context.Context, examID string) error {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	objID, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionAnswer.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam answer is removed`)
	}

	_, err = collectionAnswer.DeleteOne(ctx, filter)
	return err
}

func (e *examAnswerRepository) DeleteAllAnswerByExamID(ctx context.Context, examID string) error {
	//TODO implement me
	panic("implement me")
}
