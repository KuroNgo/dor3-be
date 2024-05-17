package quiz_question_repository

import (
	quiz_options_domain "clean-architecture/domain/quiz_options"
	quiz_question_domain "clean-architecture/domain/quiz_question"
	context2 "context"
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
	collectionQuiz     string
	collectionOptions  string
}

func NewQuizQuestionRepository(db *mongo.Database, collectionQuestion string, collectionQuiz string, collectionOptions string) quiz_question_domain.IQuizQuestionRepository {
	return &quizQuestionRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionQuiz:     collectionQuiz,
		collectionOptions:  collectionOptions,
	}
}

func (q quizQuestionRepository) FetchMany(ctx context.Context, page string) (quiz_question_domain.Response, error) {
	collectionQuestion := q.database.Collection(q.collectionQuestion)
	collectionOptions := q.database.Collection(q.collectionOptions)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return quiz_question_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// Đếm tổng số lượng tài liệu trong collection
	count, err := collectionQuestion.CountDocuments(ctx, bson.D{})
	if err != nil {
		return quiz_question_domain.Response{}, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	cursor, err := collectionQuestion.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}

	var questions []quiz_question_domain.QuizQuestionResponse
	for cursor.Next(ctx) {
		var question quiz_question_domain.QuizQuestion
		if err = cursor.Decode(&question); err != nil {
			return quiz_question_domain.Response{}, err
		}

		var option quiz_options_domain.QuizOptions
		filterOptions := bson.M{"question_id": question.ID}
		err := collectionOptions.FindOne(ctx, filterOptions).Decode(&option)
		if err != nil {
			return quiz_question_domain.Response{}, err
		}

		var questionRes quiz_question_domain.QuizQuestionResponse
		questionRes.ID = question.ID
		questionRes.QuizID = question.QuizID
		questionRes.VocabularyID = question.VocabularyID
		questionRes.Options = option
		questionRes.Content = question.Content
		questionRes.Type = question.Type
		questionRes.Level = question.Level
		questionRes.Content = question.Content
		questionRes.UpdateAt = question.UpdateAt
		questionRes.WhoUpdate = question.WhoUpdate

		questions = append(questions, questionRes)
	}
	questionsRes := quiz_question_domain.Response{
		Page:                 cal,
		QuizQuestionResponse: questions,
	}
	return questionsRes, nil
}

func (q quizQuestionRepository) FetchManyByQuizID(ctx context.Context, quizID string) (quiz_question_domain.Response, error) {
	collectionQuestion := q.database.Collection(q.collectionQuestion)
	collectionOptions := q.database.Collection(q.collectionOptions)

	idQuiz, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}

	filter := bson.M{"quiz_id": idQuiz}
	cursor, err := collectionQuestion.Find(ctx, filter)
	if err != nil {
		return quiz_question_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context2.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var questions []quiz_question_domain.QuizQuestionResponse
	for cursor.Next(ctx) {
		var question quiz_question_domain.QuizQuestion
		if err = cursor.Decode(&question); err != nil {
			return quiz_question_domain.Response{}, err
		}

		question.QuizID = idQuiz
		var option quiz_options_domain.QuizOptions
		filterOptions := bson.M{"question_id": question.ID}
		err := collectionOptions.FindOne(ctx, filterOptions).Decode(&option)
		if err != nil {
			return quiz_question_domain.Response{}, err
		}

		var questionRes quiz_question_domain.QuizQuestionResponse
		questionRes.ID = question.ID
		questionRes.QuizID = question.QuizID
		questionRes.VocabularyID = question.VocabularyID
		questionRes.Options = option
		questionRes.Content = question.Content
		questionRes.Type = question.Type
		questionRes.Level = question.Level
		questionRes.Content = question.Content
		questionRes.UpdateAt = question.UpdateAt
		questionRes.WhoUpdate = question.WhoUpdate

		questions = append(questions, questionRes)
	}

	questionsRes := quiz_question_domain.Response{
		QuizQuestionResponse: questions,
	}

	return questionsRes, nil
}

func (q quizQuestionRepository) UpdateOne(ctx context.Context, quizQuestion *quiz_question_domain.QuizQuestion) (*mongo.UpdateResult, error) {
	collection := q.database.Collection(q.collectionQuestion)

	filter := bson.D{{Key: "_id", Value: quizQuestion.ID}}
	update := bson.M{
		"$set": bson.M{
			"exam_id":    quizQuestion.QuizID,
			"content":    quizQuestion.Content,
			"level":      quizQuestion.Level,
			"update_at":  quizQuestion.UpdateAt,
			"who_update": quizQuestion.WhoUpdate,
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
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	filterExamID := bson.M{"quiz_id": quizQuestion.QuizID}
	countLessonID, err := collectionQuiz.CountDocuments(ctx, filterExamID)
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
