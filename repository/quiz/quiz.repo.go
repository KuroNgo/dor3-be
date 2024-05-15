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
	database         *mongo.Database
	collectionQuiz   string
	collectionLesson string
	collectionUnit   string
}

func NewQuizRepository(db *mongo.Database, collectionQuiz string, collectionLesson string, collectionUnit string) quiz_domain.IQuizRepository {
	return &quizRepository{
		database:         db,
		collectionQuiz:   collectionQuiz,
		collectionLesson: collectionLesson,
		collectionUnit:   collectionUnit,
	}
}

func (q *quizRepository) FetchOneByUnitID(ctx context.Context, unitID string) (quiz_domain.QuizResponse, error) {
	//e.cacheMutex.RLock()
	//cacheData, found := e.examOneCache[unitID]
	//e.cacheMutex.RUnlock()

	//if found {
	//	return cacheData, nil
	//}

	collectionQuiz := q.database.Collection(q.collectionQuiz)

	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return quiz_domain.QuizResponse{}, err
	}

	filter := bson.M{"unit_id": idUnit}
	var quiz quiz_domain.QuizResponse
	err = collectionQuiz.FindOne(ctx, filter).Decode(&quiz)
	if err != nil {
		return quiz_domain.QuizResponse{}, err
	}

	countQuestion := q.countQuestion(ctx, quiz.ID.Hex())
	quiz.CountQuestion = countQuestion

	//e.cacheMutex.Lock()
	//e.examOneCache[unitID] = exam
	//e.examCacheExpires[unitID] = time.Now().Add(5 * time.Minute)
	//e.cacheMutex.Unlock()

	return quiz, nil
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
	response := quiz_domain.Response{
		Page:        cal,
		CurrentPage: pageNumber,
	}
	return quizz, response, nil
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
		var quizRes quiz_domain.QuizResponse
		if err = cursor.Decode(&quizRes); err != nil {
			return nil, quiz_domain.Response{}, err
		}

		quiz = append(quiz, quizRes)
	}

	statisticsCh := make(chan quiz_domain.Statistics)
	go func() {
		statistics, _ := q.Statistics(ctx)
		statisticsCh <- statistics
	}()

	statistics := <-statisticsCh
	detail := quiz_domain.Response{
		Page:        cal,
		CurrentPage: pageNumber,
		Statistics:  statistics,
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

func (q *quizRepository) CreateOne(ctx context.Context, quiz *quiz_domain.Quiz) error {
	collectionQuiz := q.database.Collection(q.collectionQuiz)
	collectionUnit := q.database.Collection(q.collectionUnit)
	collectionLesson := q.database.Collection(q.collectionLesson)

	filterLesson := bson.M{"_id": quiz.LessonID}

	countL, err := collectionLesson.CountDocuments(ctx, filterLesson)
	if err != nil {
		return err
	}
	if countL > 0 {
		return errors.New("the lesson id do not exist!")
	}

	filterUnit := bson.M{"_id": quiz.UnitID}
	countN, err := collectionUnit.CountDocuments(ctx, filterUnit)
	if err != nil {
		return err
	}
	if countN > 0 {
		return errors.New("the unit id do not exist!")
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

func (q *quizRepository) countQuestion(ctx context.Context, examID string) int64 {
	collectionExamQuestion := q.database.Collection(q.collectionQuiz)

	idExam, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return 0
	}

	filter := bson.M{"quiz_id": idExam}
	count, err := collectionExamQuestion.CountDocuments(ctx, filter)
	if err != nil {
		return 0
	}

	return count
}

func (q *quizRepository) Statistics(ctx context.Context) (quiz_domain.Statistics, error) {
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	// Đếm tổng số lượng tài liệu trong collection
	count, err := collectionQuiz.CountDocuments(ctx, bson.D{})
	if err != nil {
		return quiz_domain.Statistics{}, err
	}

	statistics := quiz_domain.Statistics{
		Total: count,
	}

	return statistics, nil
}
