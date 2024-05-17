package exam_answer_repository

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	exam_question_domain "clean-architecture/domain/exam_question"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type examAnswerRepository struct {
	database           *mongo.Database
	collectionQuestion string
	collectionAnswer   string
	collectionExam     string

	answerManyCache    map[string]exam_answer_domain.ExamAnswer
	answerOneCache     map[string]exam_answer_domain.Response
	answerCacheExpires map[string]time.Time
	cacheMutex         sync.RWMutex
}

func NewExamAnswerRepository(db *mongo.Database, collectionQuestion string, collectionAnswer string, collectionExam string) exam_answer_domain.IExamAnswerRepository {
	return &examAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
		collectionExam:     collectionExam,

		answerManyCache:    make(map[string]exam_answer_domain.ExamAnswer),
		answerOneCache:     make(map[string]exam_answer_domain.Response),
		answerCacheExpires: make(map[string]time.Time),
	}
}

func (e *examAnswerRepository) FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (exam_answer_domain.Response, error) {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

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

	var answers []exam_answer_domain.ExamAnswerResponse
	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			var answer exam_answer_domain.ExamAnswer
			if err = cursor.Decode(&answer); err != nil {
				return
			}

			var question exam_question_domain.ExamQuestion
			filterQuestion := bson.M{"_id": answer.QuestionID}
			err = collectionQuestion.FindOne(ctx, filterQuestion).Decode(&question)
			if err = cursor.Decode(&answer); err != nil {
				return
			}

			var answerRes exam_answer_domain.ExamAnswerResponse
			answerRes.ID = answer.ID
			answerRes.UserID = answer.UserID
			answerRes.Question = question
			answerRes.Answer = answer.Answer
			answerRes.SubmittedAt = answer.SubmittedAt
			answerRes.IsCorrect = answer.IsCorrect

			answer.QuestionID = idQuestion
			answers = append(answers, answerRes)
		}
	}()

	internal.Wg.Wait()

	response := exam_answer_domain.Response{
		ExamAnswerResponse: answers,
	}

	return response, nil
}

func (e *examAnswerRepository) CreateOne(ctx context.Context, examAnswer *exam_answer_domain.ExamAnswer) error {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	// Kiểm tra questionID có tồn tại
	filterQuestionID := bson.M{"_id": examAnswer.QuestionID}
	countQuestionID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}
	if countQuestionID == 0 {
		return errors.New("the question ID does not exist")
	}

	// Lấy câu hỏi từ CSDL
	var examQuestion exam_question_domain.ExamQuestion
	err = collectionQuestion.FindOne(ctx, filterQuestionID).Decode(&examQuestion)
	if err != nil {
		return err
	}

	// Kiểm tra câu trả lời của thí sinh
	if examAnswer.Answer == examQuestion.CorrectAnswer {
		examAnswer.IsCorrect = 1
	}

	// Thêm câu trả lời vào CSDL
	_, err = collectionAnswer.InsertOne(ctx, examAnswer)
	if err != nil {
		return err
	}

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
	collectionAnswer := e.database.Collection(e.collectionAnswer)

	objID, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return err
	}

	filter := bson.M{"exam_id": objID}
	count, err := collectionAnswer.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam answer is removed`)
	}

	_, err = collectionAnswer.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
