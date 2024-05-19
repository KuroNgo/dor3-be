package quiz_question_repository

import (
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
}

func NewQuizQuestionRepository(db *mongo.Database, collectionQuestion string, collectionQuiz string) quiz_question_domain.IQuizQuestionRepository {
	return &quizQuestionRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionQuiz:     collectionQuiz,
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

	var questions []quiz_question_domain.QuizQuestion
	for cursor.Next(ctx) {
		var question quiz_question_domain.QuizQuestion
		if err = cursor.Decode(&question); err != nil {
			return quiz_question_domain.Response{}, err
		}

		questions = append(questions, question)
	}
	questionsRes := quiz_question_domain.Response{
		Page:         cal,
		QuizQuestion: questions,
	}
	return questionsRes, nil
}

func (q quizQuestionRepository) FetchOneByQuizID(ctx context.Context, quizID string) (quiz_question_domain.QuizQuestion, error) {
	collectionQuestion := q.database.Collection(q.collectionQuestion)

	idQuiz, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return quiz_question_domain.QuizQuestion{}, err
	}

	var quizQuestion quiz_question_domain.QuizQuestion
	filter := bson.M{"quiz_id": idQuiz}
	err = collectionQuestion.FindOne(ctx, filter).Decode(&quizQuestion)
	if err != nil {
		return quiz_question_domain.QuizQuestion{}, err
	}

	return quizQuestion, nil
}

func (q quizQuestionRepository) FetchByID(ctx context2.Context, id string) (quiz_question_domain.QuizQuestion, error) {
	collectionQuestion := q.database.Collection(q.collectionQuestion)

	idQuestion, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return quiz_question_domain.QuizQuestion{}, err
	}

	var quizQuestion quiz_question_domain.QuizQuestion
	filter := bson.M{"_id": idQuestion}
	err = collectionQuestion.FindOne(ctx, filter).Decode(&quizQuestion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return quiz_question_domain.QuizQuestion{}, errors.New("quiz question not found")
		}
		return quiz_question_domain.QuizQuestion{}, err
	}

	return quizQuestion, nil
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
	defer func(cursor *mongo.Cursor, ctx context2.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

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

	filterQuizID := bson.M{"quiz_id": quizQuestion.QuizID}
	countQuiz, err := collectionQuiz.CountDocuments(ctx, filterQuizID)
	if err != nil {
		return err
	}

	if countQuiz == 0 {
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
