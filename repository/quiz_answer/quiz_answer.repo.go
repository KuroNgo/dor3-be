package quiz_answer_repository

import (
	exam_question_domain "clean-architecture/domain/exam_question"
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	quiz_question_domain "clean-architecture/domain/quiz_question"
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

func NewQuizAnswerRepository(db *mongo.Database, collectionQuestion string, collectionAnswer string, collectionQuiz string) quiz_answer_domain.IQuizAnswerRepository {
	return &quizAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
		collectionQuiz:     collectionQuiz,
	}
}

func (q *quizAnswerRepository) FetchManyAnswerQuestionIDInUser(ctx context.Context, questionID string, userID string) (quiz_answer_domain.Response, error) {
	collectionAnswer := q.database.Collection(q.collectionAnswer)
	collectionQuestion := q.database.Collection(q.collectionQuestion)

	// chuyển đổi sang objectID từ kiểu string
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}

	// Lấy thông tin của câu hỏi từ questionID
	var question1 exam_question_domain.ExamQuestion
	filterQuestion := bson.M{"_id": idQuestion}
	err = collectionQuestion.FindOne(ctx, filterQuestion).Decode(&question1)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}

	// Lấy tất cả các câu hỏi thuộc bài kiểm tra đó
	filterExamID := bson.M{"exam_id": question1.ExamID}
	cursorQuestions, err := collectionQuestion.Find(ctx, filterExamID)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}

	defer func(cursorQuestions *mongo.Cursor, ctx context.Context) {
		err := cursorQuestions.Close(ctx)
		if err != nil {
			return
		}
	}(cursorQuestions, ctx)

	var questions []quiz_question_domain.QuizQuestion
	for cursorQuestions.Next(ctx) {
		var question quiz_question_domain.QuizQuestion
		if err := cursorQuestions.Decode(&question); err != nil {
			return quiz_answer_domain.Response{}, err
		}
		questions = append(questions, question)
	}

	// Lấy tất cả câu trả lời của người dùng cho bài kiểm tra đó
	filterAnswers := bson.M{"user_id": idUser, "question_id": bson.M{"$in": questionIDs(questions)}}
	cursorAnswers, err := collectionAnswer.Find(ctx, filterAnswers)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}
	defer cursorAnswers.Close(ctx)

	var answers []quiz_answer_domain.QuizAnswerResponse
	for cursorAnswers.Next(ctx) {
		var answer quiz_answer_domain.QuizAnswer
		if err := cursorAnswers.Decode(&answer); err != nil {
			return quiz_answer_domain.Response{}, err
		}

		var question quiz_question_domain.QuizQuestion
		filterQuestion = bson.M{"_id": answer.QuestionID}
		_ = collectionQuestion.FindOne(ctx, filterQuestion).Decode(&question)

		answerRes := quiz_answer_domain.QuizAnswerResponse{
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
		return quiz_answer_domain.Response{}, errors.New("số câu trả lời không bằng số câu hỏi nên chưa thể tạo kết quả")
	}

	// Tính điểm
	var score int
	for _, answer := range answers {
		if answer.IsCorrect == 1 {
			score++
		}
	}

	response := quiz_answer_domain.Response{
		QuizAnswer: answers,
		Score:      score,
	}

	return response, nil
}

// Hàm trợ giúp để lấy danh sách các ObjectID của các câu hỏi
func questionIDs(questions []quiz_question_domain.QuizQuestion) []primitive.ObjectID {
	ids := make([]primitive.ObjectID, len(questions))
	for i, question := range questions {
		ids[i] = question.ID
	}
	return ids
}

func (q *quizAnswerRepository) CreateOneInUser(ctx context.Context, quizAnswer *quiz_answer_domain.QuizAnswer) error {
	collectionAnswer := q.database.Collection(q.collectionAnswer)
	collectionQuestion := q.database.Collection(q.collectionQuestion)

	// kiểm tra questionId có tồn tại
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

func (q *quizAnswerRepository) DeleteOneInUser(ctx context.Context, quizID string) error {
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

func (q *quizAnswerRepository) DeleteAllAnswerByQuizIDInUser(ctx context.Context, quizId string) error {
	collectionAnswer := q.database.Collection(q.collectionAnswer)

	objID, err := primitive.ObjectIDFromHex(quizId)
	if err != nil {
		return err
	}

	filter := bson.M{"quiz_id": objID}
	count, err := collectionAnswer.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`quiz answer is removed`)
	}

	_, err = collectionAnswer.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
