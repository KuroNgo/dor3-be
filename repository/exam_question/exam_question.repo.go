package exam_question_repository

import (
	exam_question_domain "clean-architecture/domain/exam_question"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type examQuestionRepository struct {
	database             *mongo.Database
	collectionQuestion   string
	collectionExam       string
	collectionVocabulary string
}

func NewExamQuestionRepository(db *mongo.Database, collectionQuestion string, collectionExam string, collectionVocabulary string) exam_question_domain.IExamQuestionRepository {
	return &examQuestionRepository{
		database:             db,
		collectionQuestion:   collectionQuestion,
		collectionExam:       collectionExam,
		collectionVocabulary: collectionVocabulary,
	}
}

func (e *examQuestionRepository) FetchMany(ctx context.Context, page string) (exam_question_domain.Response, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectVocabulary := e.database.Collection(e.collectionVocabulary)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return exam_question_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	count, err := collectionQuestion.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	calCh := make(chan int64)
	go func() {
		defer close(calCh)
		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1 + 1
		}
	}()

	cursor, err := collectionQuestion.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	var questions []exam_question_domain.ExamQuestionResponse
	for cursor.Next(ctx) {
		var question exam_question_domain.ExamQuestionResponse
		if err = cursor.Decode(&question); err != nil {
			return exam_question_domain.Response{}, err
		}

		var question2 exam_question_domain.ExamQuestion
		if err = cursor.Decode(&question2); err != nil {
			return exam_question_domain.Response{}, err
		}

		vocab := make(chan vocabulary_domain.Vocabulary)
		go func() {
			defer close(vocab)
			filterVocab := bson.M{"_id": question2.VocabularyID}

			var vocabulary vocabulary_domain.Vocabulary
			err = collectVocabulary.FindOne(ctx, filterVocab).Decode(&vocabulary)
			if err != nil {
				return
			}

			vocab <- vocabulary
		}()

		vocabulary := <-vocab
		question.Vocabulary = vocabulary

		questions = append(questions, question)
	}

	cal := <-calCh

	questionsRes := exam_question_domain.Response{
		Count:                count,
		Page:                 cal,
		ExamQuestionResponse: questions,
	}
	return questionsRes, nil
}

func (e *examQuestionRepository) FetchManyByExamID(ctx context.Context, examID string, page string) (exam_question_domain.Response, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectVocabulary := e.database.Collection(e.collectionVocabulary)

	idExam, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return exam_question_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))
	filter := bson.M{"exam_id": idExam}

	cursor, err := collectionQuestion.Find(ctx, filter, findOptions)
	if err != nil {
		return exam_question_domain.Response{}, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println("Error closing cursor:", err)
		}
	}()

	var count int64
	count, err = collectionQuestion.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	var questions []exam_question_domain.ExamQuestionResponse
	for cursor.Next(ctx) {
		var question exam_question_domain.ExamQuestionResponse
		if err = cursor.Decode(&question); err != nil {
			return exam_question_domain.Response{}, err
		}

		question.ExamID = idExam
		var question2 exam_question_domain.ExamQuestion
		if err = cursor.Decode(&question2); err != nil {
			return exam_question_domain.Response{}, err
		}

		vocab := make(chan vocabulary_domain.Vocabulary)
		go func() {
			defer close(vocab)
			filterVocab := bson.M{"_id": question2.VocabularyID}

			var vocabulary vocabulary_domain.Vocabulary
			err = collectVocabulary.FindOne(ctx, filterVocab).Decode(&vocabulary)
			if err != nil {
				return
			}

			vocab <- vocabulary
		}()

		vocabulary := <-vocab
		question.Vocabulary = vocabulary

		questions = append(questions, question)
	}

	var cal int64
	calCh := make(chan int64)
	go func() {
		defer close(calCh)
		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1 + 1
		}
	}()
	cal = <-calCh

	questionsRes := exam_question_domain.Response{
		Count:                count,
		Page:                 cal,
		ExamQuestionResponse: questions,
	}

	return questionsRes, nil
}

func (e *examQuestionRepository) UpdateOne(ctx context.Context, examQuestion *exam_question_domain.ExamQuestion) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionQuestion)

	filter := bson.D{{Key: "_id", Value: examQuestion.ID}}
	update := bson.M{
		"$set": bson.M{
			"exam_id":    examQuestion.ExamID,
			"content":    examQuestion.Content,
			"level":      examQuestion.Level,
			"update_at":  examQuestion.UpdateAt,
			"who_update": examQuestion.WhoUpdate,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *examQuestionRepository) CreateOne(ctx context.Context, examQuestion *exam_question_domain.ExamQuestion) error {
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectionExam := e.database.Collection(e.collectionExam)
	collectionVocabulary := e.database.Collection(e.collectionVocabulary)

	filterExamID := bson.M{"_id": examQuestion.ExamID}
	countExamID, err := collectionExam.CountDocuments(ctx, filterExamID)
	if err != nil {
		return err
	}
	if countExamID == 0 {
		return errors.New("the examID do not exist")
	}

	filterVocabularyID := bson.M{"_id": examQuestion.VocabularyID}
	countVocabularyID, err := collectionVocabulary.CountDocuments(ctx, filterVocabularyID)
	if err != nil {
		return err
	}
	if countVocabularyID == 0 {
		return errors.New("the vocabularyID do not exist")
	}
	_, err = collectionQuestion.InsertOne(ctx, examQuestion)
	return nil
}

func (e *examQuestionRepository) DeleteOne(ctx context.Context, examID string) error {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	objID, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionQuestion.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam answer is removed`)
	}

	_, err = collectionQuestion.DeleteOne(ctx, filter)
	return err
}
