package exam_answer_repository

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	exam_question_domain "clean-architecture/domain/exam_question"
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

	// Chuyển đổi questionID và userID sang ObjectID
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	// Lấy thông tin của câu hỏi từ questionID
	var question1 exam_question_domain.ExamQuestion
	filterQuestion := bson.M{"_id": idQuestion}
	err = collectionQuestion.FindOne(ctx, filterQuestion).Decode(&question1)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	// Lấy tất cả các câu hỏi thuộc bài kiểm tra đó
	filterExamID := bson.M{"exam_id": question1.ExamID}
	cursorQuestions, err := collectionQuestion.Find(ctx, filterExamID)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	defer cursorQuestions.Close(ctx)

	var questions []exam_question_domain.ExamQuestion
	for cursorQuestions.Next(ctx) {
		var question exam_question_domain.ExamQuestion
		if err := cursorQuestions.Decode(&question); err != nil {
			return exam_answer_domain.Response{}, err
		}
		questions = append(questions, question)
	}

	// Lấy tất cả câu trả lời của người dùng cho bài kiểm tra đó
	filterAnswers := bson.M{"user_id": idUser, "question_id": bson.M{"$in": questionIDs(questions)}}
	cursorAnswers, err := collectionAnswer.Find(ctx, filterAnswers)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}
	defer cursorAnswers.Close(ctx)

	var answers []exam_answer_domain.ExamAnswerResponse
	for cursorAnswers.Next(ctx) {
		var answer exam_answer_domain.ExamAnswer
		if err := cursorAnswers.Decode(&answer); err != nil {
			return exam_answer_domain.Response{}, err
		}

		var question exam_question_domain.ExamQuestion
		filterQuestion := bson.M{"_id": answer.QuestionID}
		_ = collectionQuestion.FindOne(ctx, filterQuestion).Decode(&question)

		answerRes := exam_answer_domain.ExamAnswerResponse{
			ID:          answer.ID,
			UserID:      answer.UserID,
			Question:    question,
			Answer:      answer.Answer,
			SubmittedAt: answer.SubmittedAt,
			IsCorrect:   answer.IsCorrect,
		}
		answers = append(answers, answerRes)
	}

	// Kiểm tra số lượng câu hỏi và câu trả lời
	if len(answers) != len(questions) {
		return exam_answer_domain.Response{}, errors.New("số câu trả lời không bằng số câu hỏi nên chưa thể tạo kết quả")
	}

	// Tính điểm
	var score int
	for _, answer := range answers {
		if answer.IsCorrect == 1 {
			score++
		}
	}

	response := exam_answer_domain.Response{
		ExamAnswerResponse: answers,
		Score:              score,
	}

	return response, nil
}

// Hàm trợ giúp để lấy danh sách các ObjectID của các câu hỏi
func questionIDs(questions []exam_question_domain.ExamQuestion) []primitive.ObjectID {
	ids := make([]primitive.ObjectID, len(questions))
	for i, question := range questions {
		ids[i] = question.ID
	}
	return ids
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
