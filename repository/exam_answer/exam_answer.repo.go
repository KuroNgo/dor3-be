package exam_answer_repository

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	exam_options_domain "clean-architecture/domain/exam_options"
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
	collectionOptions  string
	collectionExam     string

	answerManyCache    map[string]exam_answer_domain.ExamAnswer
	answerOneCache     map[string]exam_answer_domain.Response
	answerCacheExpires map[string]time.Time
	cacheMutex         sync.RWMutex
}

func NewExamAnswerRepository(db *mongo.Database, collectionQuestion string, collectionOptions string, collectionAnswer string, collectionExam string) exam_answer_domain.IExamAnswerRepository {
	return &examAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
		collectionOptions:  collectionOptions,
		collectionExam:     collectionExam,

		answerManyCache:    make(map[string]exam_answer_domain.ExamAnswer),
		answerOneCache:     make(map[string]exam_answer_domain.Response),
		answerCacheExpires: make(map[string]time.Time),
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var answer exam_answer_domain.ExamAnswer
			if err = cursor.Decode(&answer); err != nil {
				return
			}

			// Gắn CourseID vào bài học
			answer.QuestionID = idQuestion
			answers = append(answers, answer)
		}
	}()

	wg.Wait()

	response := exam_answer_domain.Response{
		ExamAnswer: answers,
	}

	return response, nil
}

func (e *examAnswerRepository) CreateOne(ctx context.Context, examAnswer *exam_answer_domain.ExamAnswer) error {
	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionOptions := e.database.Collection(e.collectionOptions)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	// kiểm tra questionId có tồn tại
	filterQuestionID := bson.M{"question_id": examAnswer.QuestionID}
	countQuestionID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}
	if countQuestionID == 0 {
		return errors.New("the question ID do not exist")
	}

	// kiểm tra answer có bằng với đáp án
	var options exam_options_domain.ExamOptions
	if err := collectionOptions.FindOne(ctx, filterQuestionID).Decode(&options); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("no options found for the question ID")
		}
		return err
	}

	if examAnswer.Answer == options.CorrectAnswer {
		examAnswer.IsCorrect = 1 //đúng
	} else {
		examAnswer.IsCorrect = 0 //sai
	}

	// thêm answer vào CSDL
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
