package quiz_question_repository

import (
	quiz_question_domain "clean-architecture/domain/quiz_question"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"strconv"
)

type quizQuestionRepository struct {
	database           *mongo.Database
	collectionQuestion string
	collectionExam     string
}

func NewQuizQuestionRepository(db *mongo.Database, collectionQuestion string, collectionExam string) quiz_question_domain.IQuizQuestionRepository {
	return &quizQuestionRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionExam:     collectionExam,
	}
}

func (q quizQuestionRepository) FetchMany(ctx context.Context, page string) (quiz_question_domain.Response, error) {
	collectionQuestion := q.database.Collection(q.collectionQuestion)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return quiz_question_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	cursor, err := collectionQuestion.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}

	var questions []quiz_question_domain.QuizQuestion
	for cursor.Next(ctx) {
		var question quiz_question_domain.QuizQuestion
		if err = cursor.Decode(&question); err != nil {
			return quiz_question_domain.Response{}, err
		}

		questions = append(questions, question)
	}
	questionsRes := quiz_question_domain.Response{
		QuizQuestion: questions,
	}
	return questionsRes, nil
}

func (q quizQuestionRepository) FetchManyByQuizID(ctx context.Context, quizID string) (quiz_question_domain.Response, error) {
	collectionQuestion := q.database.Collection(q.collectionQuestion)

	idQuiz, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}

	filter := bson.M{"quiz_id": idQuiz}
	cursor, err := collectionQuestion.Find(ctx, filter)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var questions []quiz_question_domain.QuizQuestion
	for cursor.Next(ctx) {
		var question quiz_question_domain.QuizQuestion
		if err = cursor.Decode(&question); err != nil {
			return quiz_question_domain.Response{}, err
		}

		question.QuizID = idQuiz
		questions = append(questions, question)
	}

	questionsRes := quiz_question_domain.Response{
		QuizQuestion: questions,
	}

	return questionsRes, nil
}

func (q quizQuestionRepository) UpdateOne(ctx context.Context, quizQuestion *quiz_question_domain.QuizQuestion) (*mongo.UpdateResult, error) {
	collection := q.database.Collection(q.collectionQuestion)

	filter := bson.D{{Key: "_id", Value: quizQuestion.ID}}
	update := bson.M{
		"$set": bson.M{
			"exam_id":        quizQuestion.QuizID,
			"content":        quizQuestion.Content,
			"level":          quizQuestion.Level,
			"filename":       quizQuestion.Filename,
			"audio_duration": quizQuestion.AudioDuration,
			"update_at":      quizQuestion.UpdateAt,
			"who_update":     quizQuestion.WhoUpdate,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (q quizQuestionRepository) CreateOne(ctx context.Context, quizQuestion *quiz_question_domain.QuizQuestion) error {
	collectionQuestion := q.database.Collection(q.collectionQuestion)
	collectionExam := q.database.Collection(q.collectionExam)

	filterExamID := bson.M{"quiz_id": quizQuestion.QuizID}
	countLessonID, err := collectionExam.CountDocuments(ctx, filterExamID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the quizID do not exist")
	}

	_, err = collectionQuestion.InsertOne(ctx, quizQuestion)
	return nil
}

func (q quizQuestionRepository) DeleteOne(ctx context.Context, quizID string) error {
	collectionQuestion := q.database.Collection(q.collectionQuestion)

	objID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionQuestion.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`quiz answer is removed`)
	}

	_, err = collectionQuestion.DeleteOne(ctx, filter)
	return err
}
