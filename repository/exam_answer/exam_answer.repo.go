package exam_answer_repository

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	exam_question_domain "clean-architecture/domain/exam_question"
	"clean-architecture/internal"
	"clean-architecture/internal/cache/memory"
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
}

func NewExamAnswerRepository(db *mongo.Database, collectionQuestion string, collectionAnswer string, collectionExam string) exam_answer_domain.IExamAnswerRepository {
	return &examAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
		collectionExam:     collectionExam,
	}
}

var (
	examAnswersCache = memory.NewTTL[string, exam_answer_domain.Response]()

	wg           sync.WaitGroup
	mu           sync.Mutex
	isProcessing bool
)

const (
	cacheTTL = 5 * time.Minute
)

func (e *examAnswerRepository) FetchManyAnswerByQuestionIDInUser(ctx context.Context, questionID string, userID primitive.ObjectID) (exam_answer_domain.Response, error) {
	errCh := make(chan error, 1)
	examAnswersCh := make(chan exam_answer_domain.Response, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := examAnswersCache.Get(userID.Hex() + questionID)
		if found {
			examAnswersCh <- data
		}
	}()

	go func() {
		defer close(examAnswersCh)
		wg.Wait()
	}()

	examAnswerData := <-examAnswersCh
	if !internal.IsZeroValue(examAnswerData) {
		return examAnswerData, nil
	}

	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	// Chuyển đổi questionID và userID sang ObjectID
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
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

	defer func(cursorQuestions *mongo.Cursor, ctx context.Context) {
		err = cursorQuestions.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursorQuestions, ctx)

	var questions []exam_question_domain.ExamQuestion
	questions = make([]exam_question_domain.ExamQuestion, 0, cursorQuestions.RemainingBatchLength()) //pre-allocated slices.
	for cursorQuestions.Next(ctx) {
		var question exam_question_domain.ExamQuestion
		if err = cursorQuestions.Decode(&question); err != nil {
			return exam_answer_domain.Response{}, err
		}

		wg.Add(1)
		go func(question exam_question_domain.ExamQuestion) {
			defer wg.Done()
			questions = append(questions, question)
		}(question)
	}
	wg.Wait()

	// Lấy tất cả câu trả lời của người dùng cho bài kiểm tra đó
	filterAnswers := bson.M{"user_id": userID, "question_id": bson.M{"$in": questionIDs(questions)}}
	cursorAnswers, err := collectionAnswer.Find(ctx, filterAnswers)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}
	defer func(cursorAnswers *mongo.Cursor, ctx context.Context) {
		err = cursorAnswers.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursorAnswers, ctx)

	var answers []exam_answer_domain.ExamAnswerResponse
	answers = make([]exam_answer_domain.ExamAnswerResponse, 0, cursorAnswers.RemainingBatchLength()) //pre-allocated slices.
	for cursorAnswers.Next(ctx) {
		var answer exam_answer_domain.ExamAnswer
		if err = cursorAnswers.Decode(&answer); err != nil {
			return exam_answer_domain.Response{}, err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			var question exam_question_domain.ExamQuestion
			filterQuestion = bson.M{"_id": answer.QuestionID}
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
		}()
	}
	wg.Wait()

	// Kiểm tra số lượng câu hỏi và câu trả lời
	if len(answers) != len(questions) {
		return exam_answer_domain.Response{}, errors.New("the number of the answer is not equal to the number of the question")
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

	examAnswersCache.Set(userID.Hex()+questionID, response, cacheTTL)

	select {
	case err = <-errCh:
		return exam_answer_domain.Response{}, err
	default:
		return response, nil
	}
}

func (e *examAnswerRepository) FetchOneAnswerByQuestionIDInUser(ctx context.Context, questionID string, userID primitive.ObjectID) (exam_answer_domain.Response, error) {
	examAnswerResCh := make(chan exam_answer_domain.Response, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := examAnswersCache.Get(userID.Hex() + questionID)
		if found {
			examAnswerResCh <- data
		}
	}()

	go func() {

	}()

	collectionAnswer := e.database.Collection(e.collectionAnswer)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	filter := bson.M{"user_id": userID, "question_id": questionID}
	var examAnswer exam_answer_domain.ExamAnswer
	err := collectionAnswer.FindOne(ctx, filter).Decode(&examAnswer)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	filterQuestion := bson.M{"_id": examAnswer.QuestionID}
	var examQuestion exam_question_domain.ExamQuestion
	err = collectionQuestion.FindOne(ctx, filterQuestion).Decode(&examQuestion)
	if err != nil {
		return exam_answer_domain.Response{}, err
	}

	answer := exam_answer_domain.ExamAnswerResponse{
		ID:          examAnswer.ID,
		UserID:      userID,
		Question:    examQuestion,
		Answer:      examAnswer.Answer,
		IsCorrect:   examAnswer.IsCorrect,
		SubmittedAt: examAnswer.SubmittedAt,
	}

	var answerArr []exam_answer_domain.ExamAnswerResponse
	answerArr = append(answerArr, answer)

	examAnswerRes := exam_answer_domain.Response{
		ExamAnswerResponse: answerArr,
		Score:              answer.IsCorrect,
	}

	return examAnswerRes, nil
}

func (e *examAnswerRepository) CreateOneInUser(ctx context.Context, examAnswer *exam_answer_domain.ExamAnswer) error {
	mu.Lock()
	if isProcessing {
		mu.Unlock()
		return errors.New("another goroutine is already processing")
	}
	isProcessing = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		isProcessing = false
		mu.Unlock()
	}()

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

func (e *examAnswerRepository) DeleteOneInUser(ctx context.Context, examID string) error {
	mu.Lock()
	if isProcessing {
		mu.Unlock()
		return errors.New("another goroutine is already processing")
	}
	isProcessing = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		isProcessing = false
		mu.Unlock()
	}()

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
	if err != nil {
		return err
	}

	return nil
}

func (e *examAnswerRepository) DeleteAllAnswerByExamIDInUser(ctx context.Context, examID string) error {
	mu.Lock()
	if isProcessing {
		mu.Unlock()
		return errors.New("another goroutine is already processing")
	}
	isProcessing = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		isProcessing = false
		mu.Unlock()
	}()

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

// Hàm trợ giúp để lấy danh sách các ObjectID của các câu hỏi
func questionIDs(questions []exam_question_domain.ExamQuestion) []primitive.ObjectID {
	ids := make([]primitive.ObjectID, len(questions))
	for i, question := range questions {
		ids[i] = question.ID
	}
	return ids
}
