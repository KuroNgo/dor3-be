package exercise_answer_repository

import (
	exam_question_domain "clean-architecture/domain/exam_question"
	"clean-architecture/domain/exercise_answer"
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
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
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	// Chuyển đổi questionID và userID sang ObjectID
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return exercise_answer_domain.Response{}, err
	}

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return exercise_answer_domain.Response{}, err
	}

	// Lấy thông tin của câu hỏi từ questionID
	var question1 exam_question_domain.ExamQuestion
	filterQuestion := bson.M{"_id": idQuestion}
	err = collectionQuestion.FindOne(ctx, filterQuestion).Decode(&question1)
	if err != nil {
		return exercise_answer_domain.Response{}, err
	}

	// Lấy tất cả các câu hỏi thuộc bài kiểm tra đó
	filterExamID := bson.M{"exercise_id": question1.ExamID}
	cursorQuestions, err := collectionQuestion.Find(ctx, filterExamID)
	if err != nil {
		return exercise_answer_domain.Response{}, err
	}

	defer func(cursorQuestions *mongo.Cursor, ctx context.Context) {
		err := cursorQuestions.Close(ctx)
		if err != nil {
			return
		}
	}(cursorQuestions, ctx)

	var questions []exercise_questions_domain.ExerciseQuestion
	for cursorQuestions.Next(ctx) {
		var question exercise_questions_domain.ExerciseQuestion
		if err := cursorQuestions.Decode(&question); err != nil {
			return exercise_answer_domain.Response{}, err
		}

		questions = append(questions, question)
	}

	// Lấy tất cả câu trả lời của người dùng cho bài kiểm tra đó
	filterAnswers := bson.M{"user_id": idUser, "question_id": bson.M{"$in": questionIDs(questions)}}
	cursorAnswers, err := collectionAnswer.Find(ctx, filterAnswers)
	if err != nil {
		return exercise_answer_domain.Response{}, err
	}
	defer func(cursorAnswers *mongo.Cursor, ctx context.Context) {
		err := cursorAnswers.Close(ctx)
		if err != nil {
			return
		}
	}(cursorAnswers, ctx)

	var answers []exercise_answer_domain.ExerciseAnswerResponse
	for cursorAnswers.Next(ctx) {
		var answer exercise_answer_domain.ExerciseAnswer
		if err := cursorAnswers.Decode(&answer); err != nil {
			return exercise_answer_domain.Response{}, err
		}

		var question exercise_questions_domain.ExerciseQuestion
		filterQuestion := bson.M{"_id": answer.QuestionID}
		_ = collectionQuestion.FindOne(ctx, filterQuestion).Decode(&question)

		answerRes := exercise_answer_domain.ExerciseAnswerResponse{
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
		return exercise_answer_domain.Response{}, errors.New("số câu trả lời không bằng số câu hỏi nên chưa thể tạo kết quả")
	}

	// Tính điểm
	var score int
	for _, answer := range answers {
		if answer.IsCorrect == 1 {
			score++
		}
	}

	response := exercise_answer_domain.Response{
		ExerciseAnswer: answers,
		Score:          score,
	}

	return response, nil
}

// Hàm trợ giúp để lấy danh sách các ObjectID của các câu hỏi
func questionIDs(questions []exercise_questions_domain.ExerciseQuestion) []primitive.ObjectID {
	ids := make([]primitive.ObjectID, len(questions))
	for i, question := range questions {
		ids[i] = question.ID
	}
	return ids
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
