package quiz_repository

import (
	lesson_domain "clean-architecture/domain/lesson"
	quiz_domain "clean-architecture/domain/quiz"
	unit_domain "clean-architecture/domain/unit"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"strconv"
	"sync"
)

type quizRepository struct {
	database         *mongo.Database
	collectionQuiz   string
	collectionLesson string
	collectionUnit   string
}

func (q *quizRepository) FetchByID(ctx context.Context, id string) (quiz_domain.Quiz, error) {
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	idQuiz, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return quiz_domain.Quiz{}, err
	}

	var quiz quiz_domain.Quiz
	filter := bson.M{"_id": idQuiz}
	err = collectionQuiz.FindOne(ctx, filter).Decode(&quiz)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return quiz_domain.Quiz{}, errors.New("quiz not found")
		}
		return quiz_domain.Quiz{}, err
	}

	return quiz, nil
}

func NewQuizRepository(db *mongo.Database, collectionQuiz string, collectionLesson string, collectionUnit string) quiz_domain.IQuizRepository {
	return &quizRepository{
		database:         db,
		collectionQuiz:   collectionQuiz,
		collectionLesson: collectionLesson,
		collectionUnit:   collectionUnit,
	}
}

func (q *quizRepository) FetchOneByUnitID(ctx context.Context, unitID string) (quiz_domain.Quiz, error) {
	//e.cacheMutex.RLock()
	//cacheData, found := e.examOneCache[unitID]
	//e.cacheMutex.RUnlock()

	//if found {
	//	return cacheData, nil
	//}
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return quiz_domain.Quiz{}, err
	}

	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionQuiz.Find(ctx, filter)
	if err != nil {
		return quiz_domain.Quiz{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var questions []quiz_domain.Quiz

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var quiz quiz_domain.Quiz
			if err = cursor.Decode(&quiz); err != nil {
				return
			}

			questions = append(questions, quiz)
		}
	}()
	wg.Wait()

	// Kiểm tra nếu danh sách exercises không rỗng
	if len(questions) == 0 {
		return quiz_domain.Quiz{}, errors.New("no exercises found")
	}

	// Chọn một giá trị ngẫu nhiên từ danh sách exercises
	randomIndex := rand.Intn(len(questions))
	randomExercise := questions[randomIndex]

	return randomExercise, nil
}

func (q *quizRepository) FetchManyByUnitID(ctx context.Context, unitID string, page string) ([]quiz_domain.Quiz, quiz_domain.Response, error) {
	collectionQuiz := q.database.Collection(q.collectionQuiz)
	collectionUnit := q.database.Collection(q.collectionUnit)
	collectionLesson := q.database.Collection(q.collectionLesson)

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

	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}

	cursor, err := collectionQuiz.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, quiz_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var quizz []quiz_domain.Quiz
	for cursor.Next(ctx) {
		var quiz quiz_domain.Quiz
		if err := cursor.Decode(&quiz); err != nil {
			return nil, quiz_domain.Response{}, err
		}

		var unit unit_domain.Unit
		if err = collectionUnit.FindOne(ctx, bson.M{"_id": idUnit}).Decode(&unit); err != nil {
			return nil, quiz_domain.Response{}, err
		}

		var lesson lesson_domain.Lesson
		if err = collectionLesson.FindOne(ctx, bson.M{"_id": unit.LessonID}).Decode(&lesson); err != nil {
			return nil, quiz_domain.Response{}, err
		}

		quiz.UnitID = unit.ID
		quiz.LessonID = lesson.ID

		quizz = append(quizz, quiz)
	}

	cal := <-calCh
	response := quiz_domain.Response{
		Page:        cal,
		CurrentPage: pageNumber,
	}
	return quizz, response, nil
}

func (q *quizRepository) FetchMany(ctx context.Context, page string) ([]quiz_domain.Quiz, quiz_domain.Response, error) {
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

	var quiz []quiz_domain.Quiz

	for cursor.Next(ctx) {
		var quizRes quiz_domain.Quiz
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

	filter := bson.D{{Key: "_id", Value: quiz.Id}}
	update := bson.M{"$set": quiz}

	data, err := collectionQuiz.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
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
