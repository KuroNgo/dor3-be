package quiz_answer_repository

import (
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	quiz_options_domain "clean-architecture/domain/quiz_options"
	"clean-architecture/internal"
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
	collectionOptions  string
}

func NewQuizAnswerRepository(db *mongo.Database, collectionQuestion string, collectionAnswer string, collectionOptions string) quiz_answer_domain.IQuizAnswerRepository {
	return &quizAnswerRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionAnswer:   collectionAnswer,
		collectionOptions:  collectionOptions,
	}
}

func (q *quizAnswerRepository) FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (quiz_answer_domain.Response, error) {
	collectionAnswer := q.database.Collection(q.collectionAnswer)

	// chuyển đổi sang objectID từ kiểu string
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return quiz_answer_domain.Response{}, err
	}

	// tìm kiếm
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
	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			var answer quiz_answer_domain.QuizAnswer
			if err = cursor.Decode(&answer); err != nil {
				return
			}

			// Gắn CourseID vào bài học
			answer.QuestionID = idQuestion
			answers = append(answers, answer)
		}
	}()
	internal.Wg.Wait()

	response := quiz_answer_domain.Response{
		QuizAnswer: answers,
	}

	return response, nil
}

func (q *quizAnswerRepository) CreateOne(ctx context.Context, quizAnswer *quiz_answer_domain.QuizAnswer) error {
	collectionAnswer := q.database.Collection(q.collectionAnswer)
	collectionQuestion := q.database.Collection(q.collectionQuestion)
	collectionOptions := q.database.Collection(q.collectionOptions)

	// kiểm tra questionId có tồn tại
	filterQuestionID := bson.M{"question_id": quizAnswer.QuestionID}
	countQuestionID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}
	if countQuestionID == 0 {
		return errors.New("the question ID do not exist")
	}

	// kiểm tra answer có bằng với đáp án
	var options quiz_options_domain.QuizOptions
	if err := collectionOptions.FindOne(ctx, filterQuestionID).Decode(&options); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("no options found for the question ID")
		}
		return err
	}

	if quizAnswer.Answer == options.CorrectAnswer {
		quizAnswer.IsCorrect = 1 //đúng
	} else {
		quizAnswer.IsCorrect = 0 //sai
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

func (q *quizAnswerRepository) DeleteAllAnswerByQuizID(ctx context.Context, quizId string) error {
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
