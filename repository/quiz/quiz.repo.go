package quiz_repository

import (
	quiz_domain "clean-architecture/domain/quiz"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type quizRepository struct {
	database             *mongo.Database
	collectionLesson     string
	collectionUnit       string
	collectionVocabulary string
	collectionQuiz       string
}

func (q *quizRepository) FetchOneByUnitID(ctx context.Context, unitID string) (quiz_domain.QuizResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewQuizRepository(db *mongo.Database, collectionQuiz string, collectionLesson string, collectionUnit string, collectionVocabulary string) quiz_domain.IQuizRepository {
	return &quizRepository{
		database:             db,
		collectionQuiz:       collectionQuiz,
		collectionLesson:     collectionLesson,
		collectionUnit:       collectionUnit,
		collectionVocabulary: collectionVocabulary,
	}
}
func (q *quizRepository) FetchManyByLessonID(ctx context.Context, unitID string, page string) ([]quiz_domain.QuizResponse, quiz_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (q *quizRepository) UpdateCompleted(ctx context.Context, quiz *quiz_domain.Quiz) error {
	collection := q.database.Collection(q.collectionUnit)

	filter := bson.D{{Key: "_id", Value: quiz.ID}}
	update := bson.M{"$set": bson.M{
		"is_complete": quiz.IsComplete,
		"who_updates": quiz.WhoUpdates,
	}}

	_, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return err
	}

	return nil
}

func (q *quizRepository) FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]quiz_domain.QuizResponse, quiz_domain.Response, error) {
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, quiz_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.
		Find().
		SetLimit(int64(perPage)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{"level", 1}})

	calCh := make(chan int64)
	countUnitCh := make(chan int64)

	go func() {
		defer close(calCh)
		defer close(countUnitCh)
		count, err := collectionQuiz.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1
		}
	}()

	idLesson2, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}

	filter := bson.M{"lesson_id": idLesson2}
	cursor, err := collectionQuiz.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var quizz []quiz_domain.QuizResponse
	for cursor.Next(ctx) {
		var quiz quiz_domain.QuizResponse
		if err := cursor.Decode(&quiz); err != nil {
			return nil, quiz_domain.Response{}, err
		}

		// Gắn LessonID vào đơn vị
		quiz.LessonID = idLesson2

		quizz = append(quizz, quiz)
	}

	cal := <-calCh
	totalUnit := <-countUnitCh
	response := quiz_domain.Response{
		Total:       totalUnit,
		Page:        cal,
		CurrentPage: pageNumber,
	}
	return quizz, response, nil
}

func (q *quizRepository) FetchTenQuizButEnoughAllSkill(ctx context.Context) ([]quiz_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (q *quizRepository) FetchMany(ctx context.Context, page string) ([]quiz_domain.QuizResponse, quiz_domain.Response, error) {
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, quiz_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// Đếm tổng số lượng tài liệu trong collection
	count, err := collectionQuiz.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	cursor, err := collectionQuiz.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var quiz []quiz_domain.QuizResponse

	for cursor.Next(ctx) {
		var q quiz_domain.QuizResponse
		if err := cursor.Decode(&q); err != nil {
			return nil, quiz_domain.Response{}, err
		}
		quiz = append(quiz, q)
	}

	detail := quiz_domain.Response{
		Total:       count,
		Page:        cal,
		CurrentPage: pageNumber,
	}

	return quiz, detail, nil
}

func (q *quizRepository) UpdateOne(ctx context.Context, quiz *quiz_domain.Quiz) (*mongo.UpdateResult, error) {
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	filter := bson.D{{Key: "_id", Value: quiz.ID}}
	update := bson.M{"$set": quiz}

	data, err := collectionQuiz.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (q *quizRepository) CreateOne(ctx context.Context, quiz *quiz_domain.Quiz) error {
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	filter := bson.M{"lesson_id": quiz.LessonID}

	// check exists with Count Documents
	count, err := collectionQuiz.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the question did exist")
	}

	_, err = collectionQuiz.InsertOne(ctx, quiz)
	return err
}

func (q *quizRepository) DeleteOne(ctx context.Context, quizID string) error {
	collectionQuiz := q.database.Collection(q.collectionQuiz)
	objID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}
	count, err := collectionQuiz.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the quiz is removed`)
	}
	_, err = collectionQuiz.DeleteOne(ctx, filter)
	return err
}
