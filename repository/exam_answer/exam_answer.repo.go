package exam_answer_repository

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type examAnswerRepository struct {
	database           mongo.Database
	collectionQuestion string
	collectionAnswer   string
}

func NewExamAnswerRepository(db mongo.Database, collectionQuestion string, collectionAnswer string) exam_answer_domain.IExamAnswerRepository {
	return &examAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
	}
}

func (e *examAnswerRepository) FetchManyByQuestionID(ctx context.Context, questionID string) (exam_answer_domain.Response, error) {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	filter := bson.M{"question_id": idQuestion}
	cursor, err := collectionAnswer.Find(ctx, filter)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}
	defer cursor.Close(ctx)

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

func (e *examAnswerRepository) UpdateOne(ctx context.Context, examAnswerID string, examAnswer exam_answer_domain.ExamAnswer) error {
	collection := e.database.Collection(e.collectionAnswer)
	doc, err := internal.ToDoc(examAnswer)
	objID, err := primitive.ObjectIDFromHex(examAnswerID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (e *examAnswerRepository) CreateOne(ctx context.Context, examAnswer *exam_answer_domain.ExamAnswer) error {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	filterQuestionID := bson.M{"question_id": examAnswer.QuestionID}
	countLessonID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the question ID do not exist")
	}

	_, err = collectionAnswer.InsertOne(ctx, examAnswer)
	return nil
}

func (e *examAnswerRepository) DeleteOne(ctx context.Context, examID string) error {
	collectionExam := e.database.Collection(e.collectionAnswer)
	objID, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionExam.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam answer is removed`)
	}

	_, err = collectionExam.DeleteOne(ctx, filter)
	return err
}
